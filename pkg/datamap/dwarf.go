package datamap

import (
	"debug/dwarf"
	"debug/elf"
	"errors"
	"fmt"
	"io"

	"github.com/parca-dev/runtime-data/pkg/symbols"
)

func (dataMap *DataMap) ReadFromDWARF(ef *elf.File) error {
	dwarfData, err := ef.DWARF()
	if err != nil {
		return fmt.Errorf("failed to read DWARF info: %w", err)
	}

	p := processor{
		ef:        ef,
		dwarfData: dwarfData,
	}
	// TODO(kakkoyun): This is a very naive implementation.
	// We should optimize this for performance.
	for _, rn := range dataMap.Routes {
		entries, err := p.findCompositeTypeEntries(rn.Type)
		if err != nil {
			continue
		}

		entry, err := p.findActionableEntry(entries)
		if err != nil {
			return fmt.Errorf("failed to find composite type (%s): %w", rn.Type, err)
		}

		typ, err := dwarfData.Type(entry.Offset)
		if err != nil {
			return fmt.Errorf("failed to get type: %w", err)
		}

		if err := p.process(rn, entry, typ, 0); err != nil {
			return fmt.Errorf("failed to process: %w", err)
		}
	}
	return nil
}

func (p *processor) process(rn *RouteNode, entry *dwarf.Entry, typ dwarf.Type, offset int64) error {
	if rn == nil {
		return nil
	}

	st := typ.(*dwarf.StructType)
	if rn.IsLeaf() {
		if err := p.extract(rn, entry, st, offset); err != nil {
			return fmt.Errorf("failed to extract: %w", err)
		}
		return nil
	}

	fields := map[string]*dwarf.StructField{}
	for _, f := range st.Field {
		fields[f.Name] = f
	}

	field, ok := fields[rn.Next.Type]
	if !ok {
		return fmt.Errorf("field %s not found in %s", rn.Next.Type, rn.Type)
	}

	fieldEntry, err := p.findFieldEntry(entry, field.Name)
	if err != nil {
		return fmt.Errorf("failed to find field (%s) entry: %w", field.Name, err)
	}

	return p.process(rn.Next, fieldEntry, field.Type, offset+field.ByteOffset)
}

func (p *processor) extract(rn *RouteNode, entry *dwarf.Entry, st *dwarf.StructType, offset int64) error {
	fields := map[string]*dwarf.StructField{}
	for _, f := range st.Field {
		fields[f.Name] = f
	}
	for _, ex := range rn.Extractors {
		if ex.Op == OpSizeOf {
			if ex.Source == rn.Type {
				if err := ex.Set(int64(st.Size())); err != nil {
					return fmt.Errorf("failed to set size: %w", err)
				}
				continue
			}

			field, ok := fields[ex.Source]
			if !ok {
				return fmt.Errorf("field %s not found in %s", ex.Source, rn.Type)
			}
			if err := ex.Set(int64(field.Type.Size())); err != nil {
				return fmt.Errorf("failed to set size: %w", err)
			}
		}

		if ex.Op == OpOffsetOf {
			if ex.Static {
				_ = p.extractStatic(entry, st, ex)
				continue
			}

			field, ok := fields[ex.Source]
			if !ok {
				return fmt.Errorf("field %s not found in %s", ex.Source, rn.Type)
			}

			if err := ex.Set(int64(offset + field.ByteOffset)); err != nil {
				return fmt.Errorf("failed to set offset: %w", err)
			}
		}
	}
	return nil
}

func (p *processor) extractStatic(entry *dwarf.Entry, st *dwarf.StructType, ex *Extractor) error {
	name := ex.Source
	fieldEntry, err := p.findFieldEntry(entry, name)
	if err != nil {
		return fmt.Errorf("failed to find field (%s.%s) entry: %w", st.StructName, name, err)
	}
	attributes := attrs(fieldEntry)
	linkageNameAttr, ok := attributes[dwarf.AttrLinkageName]
	if !ok {
		return fmt.Errorf("no linkage name attribute for %s", name)
	}

	linkageName := linkageNameAttr.(string)
	sym, err := symbols.FindSymbol(p.ef, linkageName)
	if err != nil {
		return fmt.Errorf("failed to find symbol (%s): %w", linkageName, err)
	}
	if err := ex.Set(int64(sym.Value)); err != nil {
		return fmt.Errorf("failed to set offset of (%s.%s): %w", st.StructName, name, err)
	}
	return nil
}

func isCompositeType(entry *dwarf.Entry) bool {
	return entry.Tag == dwarf.TagStructType || entry.Tag == dwarf.TagClassType
}

func isDeclaration(entry *dwarf.Entry) bool {
	attributes := attrs(entry)
	_, ok := attributes[dwarf.AttrDeclaration]
	return ok
}

type processor struct {
	ef        *elf.File
	dwarfData *dwarf.Data
}

// findCompositeTypeEntries finds the entries with the given name in the DWARF data.
func (p *processor) findCompositeTypeEntries(name string) ([]*dwarf.Entry, error) {
	entries := []*dwarf.Entry{}
	entryReader := p.dwarfData.Reader()
	for {
		entry, err := entryReader.Next()
		if err == io.EOF || entry == nil {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("unexpected error while reading DWARF data: %w", err)
		}

		if !isCompositeType(entry) && entry.Tag != dwarf.TagTypedef {
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

// findActionableEntry finds the composite type in the given entries.
func (p *processor) findActionableEntry(entries []*dwarf.Entry) (*dwarf.Entry, error) {
	for _, entry := range entries {
		if isCompositeType(entry) {
			// fast path.
			if !entry.Children {
				continue
			}
			if isDeclaration(entry) {
				continue
			}
			return entry, nil
		}

		typeEntry, err := typeOf(p.dwarfData, entry)
		if err != nil {
			return nil, fmt.Errorf("failed to get type: %w", err)
		}

		if !isCompositeType(typeEntry) {
			continue
		}

		if !typeEntry.Children {
			continue
		}

		if isDeclaration(typeEntry) {
			continue
		}

		return typeEntry, nil
	}
	return nil, errors.New("no composite(struct|class) type found")
}

func (p *processor) findFieldEntry(entry *dwarf.Entry, name string) (*dwarf.Entry, error) {
	entryReader := p.dwarfData.Reader()
	entryReader.Seek(entry.Offset)
	for {
		entry, err := entryReader.Next()
		if err != nil {
			return nil, err
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

		return entry, nil
	}
	return nil, errors.New("not found")
}

// Helpers:

func attrs(entry *dwarf.Entry) map[dwarf.Attr]any {
	attrs := map[dwarf.Attr]any{}
	for f := range entry.Field {
		if _, ok := attrs[entry.Field[f].Attr]; ok {
			panic(fmt.Sprintf("duplicate attribute: %s", entry.Field[f].Attr))
		}
		attrs[entry.Field[f].Attr] = entry.Field[f].Val
	}
	return attrs
}

func nameAttr(attrs map[dwarf.Attr]any) string {
	nameAttr, ok := attrs[dwarf.AttrName]
	if !ok {
		return ""
	}
	return nameAttr.(string)
}

func typeOf(dwarfData *dwarf.Data, entry *dwarf.Entry) (*dwarf.Entry, error) {
	attrs := attrs(entry)
	typeAttr, ok := attrs[dwarf.AttrType]
	if !ok {
		return nil, fmt.Errorf("no type attribute found for (%s)", nameAttr(attrs))
	}
	typeReader := dwarfData.Reader()
	typeReader.Seek(typeAttr.(dwarf.Offset))
	typeEntry, err := typeReader.Next()
	if err != nil {
		return nil, fmt.Errorf("unexpected error while reading DWARF data: %w", err)
	}
	return typeEntry, nil
}
