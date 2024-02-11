package datamap

import (
	"debug/dwarf"
	"errors"
	"fmt"
	"io"
)

// ReadFromDWARF reads the DWARF data and extracts the offsets of the definitions.
func (dataMap *DataMap) ReadFromDWARF(dwarfData *dwarf.Data) error {
	// TODO: Optimize finding the entries and struct types.
	// - Wait until we are sure about the correctness of the implementation.
	for _, rn := range dataMap.Routes {
		entries, err := findEntries(dwarfData, rn.Type)
		if err != nil {
			continue
		}

		entry, err := findStructType(dwarfData, entries)
		if err != nil {
			return fmt.Errorf("failed to find struct type (%s): %w", rn.Type, err)
		}

		if err := process(dwarfData, entry, rn); err != nil {
			if errors.Is(err, errNoSize) {
				continue
			}
			return fmt.Errorf("failed to extract: %w", err)
		}
	}
	return nil
}

// findEntries finds the entries with the given name in the DWARF data.
func findEntries(dwarfData *dwarf.Data, name string) ([]*dwarf.Entry, error) {
	entries := []*dwarf.Entry{}
	entryReader := dwarfData.Reader()
	for {
		entry, err := entryReader.Next()
		if err == io.EOF || entry == nil {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("unexpected error while reading DWARF data: %w", err)
		}

		if entry.Tag != dwarf.TagStructType && entry.Tag != dwarf.TagTypedef {
			continue
		}

		attributes := attrs(entry)
		val, ok := attributes[dwarf.AttrName]
		if !ok {
			continue
		}

		if val.(string) != name {
			continue
		}

		entries = append(entries, entry)
	}
	return entries, nil
}

// findStructType finds the struct type in the given entries.
func findStructType(dwarfData *dwarf.Data, entries []*dwarf.Entry) (*dwarf.Entry, error) {
	for _, entry := range entries {
		if entry.Tag == dwarf.TagStructType {
			if !entry.Children {
				continue
			}
			return entry, nil
		}

		typeEntry, err := typeOf(dwarfData, entry)
		if err != nil {
			return nil, fmt.Errorf("failed to get type: %w", err)
		}

		if typeEntry.Tag != dwarf.TagStructType {
			continue
		}

		if !typeEntry.Children {
			continue
		}

		return typeEntry, nil
	}
	return nil, errors.New("no struct type found")
}

// process processes the given node in the route using the DWARF data.
func process(dwarfData *dwarf.Data, entry *dwarf.Entry, rn *RouteNode) error {
	if rn.IsLeaf() {
		// We are at the end of the path,
		// and we have the type we need to extract the data.
		return processLeaf(dwarfData, entry, rn, 0)
	}

	// offset, err := offsetOf(dwarfData, entry, rn.Next.Type)
	// if err != nil {
	// 	return fmt.Errorf("failed to get offset of field (%s): %w", rn.Next.Type, err)
	// }
	return processNested(dwarfData, entry, rn.Next, 0)
}

var errNotFound = errors.New("not found")

// processNested finds the nested struct in the given struct.
func processNested(dwarfData *dwarf.Data, entry *dwarf.Entry, rn *RouteNode, offset int64) error {
	// entry is the struct we are inside.
	// rn is the struct we are looking for inside the entry.
	name := rn.Type
	entryReader := dwarfData.Reader()
	entryReader.Seek(entry.Offset)
	for {
		entry, err := entryReader.Next()
		if err != nil {
			return err
		}

		if err == io.EOF || entry == nil {
			break
		}

		if entry.Tag == 0 {
			// End of children.
			break
		}

		attributes := attrs(entry)
		nameAttr, ok := attributes[dwarf.AttrName]
		if !ok {
			continue
		}

		fieldName := nameAttr.(string)
		if len(fieldName) == 0 {
			continue
		}

		if fieldName != name {
			continue
		}

		// Found. Now jump to the type and process it.

		typeAttr, ok := attributes[dwarf.AttrType]
		if !ok {
			continue
		}

		typeReader := dwarfData.Reader()
		typeReader.Seek(typeAttr.(dwarf.Offset))
		typeEntry, err := typeReader.Next()
		if err != nil {
			return err
		}

		if !typeEntry.Children {
			continue
		}

		offsetAttr, ok := attributes[dwarf.AttrDataMemberLoc]
		if !ok {
			return fmt.Errorf("no offset attribute for %s", name)
		}
		offset = offset + offsetAttr.(int64)

		if rn.IsLeaf() {
			// We are at the end of the path,
			// and we have the type we need to extract the data.
			return processLeaf(dwarfData, typeEntry, rn, offset)
		} else {
			return processNested(dwarfData, typeEntry, rn.Next, offset)
		}
	}

	return errNotFound
}

var errNoSize = errors.New("no size")

// processLeaf handles the extraction of the offset or size of the struct fields at the leaf level.
func processLeaf(dwarfData *dwarf.Data, entry *dwarf.Entry, rn *RouteNode, offset int64) error {
	fieldSourcesQueriedFor := map[string]*Extractor{}
	typeSourcesQueriedFor := map[string]*Extractor{}
	for _, f := range rn.Extractors {
		if f.Source == rn.Type {
			typeSourcesQueriedFor[f.Source] = f
		} else {
			fieldSourcesQueriedFor[f.Source] = f
		}
	}

	// Process type extractors.
	for _, ext := range typeSourcesQueriedFor {
		if ext.Op != OpSizeOf {
			// We only support sizeof for types.
			continue
		}

		size, err := sizeOf(dwarfData, entry)
		if err != nil {
			return fmt.Errorf("failed to get size: %w", err)
		}

		if err := ext.Set(size); err != nil {
			return fmt.Errorf("failed to set size: %w", err)
		}
	}

	attributes := attrs(entry)
	sizeAttr, ok := attributes[dwarf.AttrByteSize]
	if !ok {
		return errNoSize
	}
	size := sizeAttr.(int64)
	if size == 0 && !entry.Children {
		// Skip children if the size is 0.
		return nil
	}

	// Read the fields of the struct.
	entryReader := dwarfData.Reader()
	entryReader.Seek(entry.Offset)
	for {
		entry, err := entryReader.Next()
		if err != nil {
			return err
		}

		if err == io.EOF || entry == nil {
			break
		}

		if entry.Tag == 0 {
			// End of children.
			break
		}

		attributes := attrs(entry)
		nameAttr, ok := attributes[dwarf.AttrName]
		if !ok {
			continue
		}
		fieldName := nameAttr.(string)
		if len(fieldName) == 0 {
			continue
		}
		field, ok := fieldSourcesQueriedFor[fieldName]
		if !ok {
			continue
		}

		switch field.Op {
		case OpOffsetOf:
			offsetAttr, ok := attributes[dwarf.AttrDataMemberLoc]
			if !ok {
				continue
			}
			offset := offsetAttr.(int64) + offset
			if err := field.Set(offset); err != nil {
				return fmt.Errorf("failed to set offset: %w", err)
			}
		case OpSizeOf:
			size, err := sizeOf(dwarfData, entry)
			if err != nil {
				return fmt.Errorf("failed to get size: %w", err)
			}
			if err := field.Set(size); err != nil {
				return fmt.Errorf("failed to set size: %w", err)
			}
		}
	}
	return nil
}

// Helpers:

func attrs(entry *dwarf.Entry) map[dwarf.Attr]any {
	attrs := map[dwarf.Attr]any{}
	for f := range entry.Field {
		if _, ok := attrs[entry.Field[f].Attr]; ok {
			panic("duplicate attr")
		}
		attrs[entry.Field[f].Attr] = entry.Field[f].Val
	}
	return attrs
}

func typeOf(dwarfData *dwarf.Data, entry *dwarf.Entry) (*dwarf.Entry, error) {
	attrs := attrs(entry)
	typeAttr, ok := attrs[dwarf.AttrType]
	if !ok {
		return nil, errors.New("no type attribute")
	}
	typeReader := dwarfData.Reader()
	typeReader.Seek(typeAttr.(dwarf.Offset))
	typeEntry, err := typeReader.Next()
	if err != nil {
		return nil, fmt.Errorf("unexpected error while reading DWARF data: %w", err)
	}
	return typeEntry, nil
}

func sizeOf(dwarfData *dwarf.Data, entry *dwarf.Entry) (int64, error) {
	attrs := attrs(entry)
	sizeAttr, ok := attrs[dwarf.AttrByteSize]
	if ok {
		return sizeAttr.(int64), nil
	}

	typeEntry, err := typeOf(dwarfData, entry)
	if err != nil {
		return 0, fmt.Errorf("failed to get type: %w", err)
	}

	return sizeOf(dwarfData, typeEntry)
}
