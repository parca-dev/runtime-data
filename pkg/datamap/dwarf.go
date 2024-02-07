package datamap

import (
	"debug/dwarf"
	"errors"
	"fmt"
	"io"
)

// ReadFromDWARF reads the DWARF data and extracts the offsets of the definitions.
func (dataMap *DataMap) ReadFromDWARF(dwarfData *dwarf.Data) error {
	namesQueriedFor := make(map[string]*Struct, len(dataMap.Structs))
	for _, sdm := range dataMap.Structs {
		namesQueriedFor[sdm.StructName] = sdm
	}
	if len(namesQueriedFor) == 0 {
		return errors.New("no struct names provided")
	}

	alreadyExtracted := map[string]struct{}{}
	entryReader := dwarfData.Reader()
	typeReader := dwarfData.Reader()
	for {
		entry, err := entryReader.Next()
		if err == io.EOF || entry == nil {
			break
		}

		if err != nil {
			return fmt.Errorf("unexpected error while reading DWARF data: %w", err)
		}

		if entry.Tag != dwarf.TagStructType && entry.Tag != dwarf.TagTypedef {
			continue
		}

		attributes := attrs(entry)
		switch entry.Tag {
		case dwarf.TagTypedef:
			val, ok := attributes[dwarf.AttrName]
			if !ok {
				continue
			}

			typeName := val.(string)
			if len(typeName) == 0 {
				continue
			}
			sm, ok := namesQueriedFor[typeName]
			if !ok {
				continue
			}
			if _, ok := alreadyExtracted[typeName]; ok {
				continue
			}

			typeAttr := attributes[dwarf.AttrType]
			if typeAttr == nil {
				continue
			}

			typeReader.Seek(typeAttr.(dwarf.Offset))
			typeEntry, err := typeReader.Next()
			if err != nil {
				return fmt.Errorf("unexpected error while reading DWARF data: %w", err)
			}

			if typeEntry.Tag != dwarf.TagStructType {
				return fmt.Errorf("unexpected tag, only structs are supported: %v", typeEntry.Tag)
			}

			if !typeEntry.Children {
				continue
			}

			if err := extractFromStructEntry(typeReader, typeEntry, sm); err != nil {
				if errors.Is(err, errNoSize) {
					continue
				}
				return fmt.Errorf("failed to extract field offsets: %w", err)
			}

			alreadyExtracted[typeName] = struct{}{}
		case dwarf.TagStructType:
			val, ok := attributes[dwarf.AttrName]
			if !ok {
				continue
			}

			structName := val.(string)
			if len(structName) == 0 {
				continue
			}
			if _, ok := namesQueriedFor[structName]; !ok {
				continue
			}
			if _, ok := alreadyExtracted[structName]; ok {
				continue
			}

			sm := namesQueriedFor[structName]
			if sm == nil {
				return fmt.Errorf("struct %s not found", structName)
			}

			if err := extractFromStructEntry(entryReader, entry, sm); err != nil {
				if errors.Is(err, errNoSize) {
					continue
				}
				return fmt.Errorf("failed to extract field offsets: %w", err)
			}

			alreadyExtracted[structName] = struct{}{}
		default:
		}
	}
	return nil
}

var errNoSize = errors.New("no size")

// extractFromStructEntry handles the extraction of the offset or size of the struct.
func extractFromStructEntry(entryReader *dwarf.Reader, entry *dwarf.Entry, sm *Struct) error {
	structName := sm.StructName
	attributes := attrs(entry)

	sizeAttr, ok := attributes[dwarf.AttrByteSize]
	if !ok {
		return errNoSize
	}
	size := sizeAttr.(int64)
	if size == 0 && entry.Children {
		entryReader.SkipChildren()
		return nil
	}

	switch sm.Op {
	case OpOffset:
		if err := extractStructFieldOffsets(entryReader, entry, sm); err != nil {
			return fmt.Errorf("failed to extract field offsets: %w", err)
		}
	case OpSize:
		if sm.Value.CanSet() {
			sm.Value.SetInt(size)
		} else {
			return fmt.Errorf("size of struct %s is not settable", structName)
		}
	}
	return nil
}

// extractStructFieldOffsets extracts the field offsets from the DWARF struct type entry.
func extractStructFieldOffsets(entryReader *dwarf.Reader, entry *dwarf.Entry, sdm *Struct) error {
	namesQueriedFor := make(map[string]*Field, len(sdm.Fields))
	for _, f := range sdm.Fields {
		namesQueriedFor[f.Name] = f
	}

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

		if entry.Tag != dwarf.TagMember {
			panic("unexpected tag")
		}

		attributes := attrs(entry)
		offsetAttr, ok := attributes[dwarf.AttrName]
		if !ok {
			continue
		}
		fieldName := offsetAttr.(string)
		if len(fieldName) == 0 {
			continue
		}
		field, ok := namesQueriedFor[fieldName]
		if !ok {
			continue
		}

		offsetAttr, ok = attributes[dwarf.AttrDataMemberLoc]
		if !ok {
			continue
		}
		offset := offsetAttr.(int64)
		if field.Value.CanSet() {
			field.Value.SetInt(offset)
		} else {
			return fmt.Errorf("field %s is not settable", fieldName)
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
