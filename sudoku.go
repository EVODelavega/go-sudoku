package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

/**
 * Solve sudoku puzzles by feeding it the raw input (either file or as CLI arg)
 * Use -h[elp] or no arguments to see expected flags
 * Expected format: \d{81}, 1-9 for known fields, 0 for unknown values
 *
 * Output: The given (unsolved sudoku), followed by the solution
 */

type Grid struct {
	Fields  [81]Field
	HRows   [9]Row
	VRows   [9]Row
	Groups  [9]Group
	Options map[int]int
}

type Row struct {
	Fields  [9]*Field
	Options map[int]bool
}

type Group struct {
	Fields  [9]*Field
	Options map[int]bool
}

type Field struct {
	Value   int
	HRow    *Row
	VRow    *Row
	Group   *Group
	Options map[int]bool
	Grid    *Grid
}

func getEmptyOptionsMap() map[int]bool {
	return map[int]bool{
		1: true,
		2: true,
		3: true,
		4: true,
		5: true,
		6: true,
		7: true,
		8: true,
		9: true,
	}
}

//@todo -> fix error messages
func (gr *Group) AppendField(f *Field) error {
	for i, sf := range gr.Fields {
		if sf == nil {
			gr.Fields[i] = f
			return nil
		} else if sf == f {
			return fmt.Errorf("Field already added to group")
		}
	}
	return fmt.Errorf("Group already 'full'")
}

func (g *Grid) GetField(x, y int) (*Field, error) {
	offset := x + 9*y
	if offset > 80 {
		return nil, fmt.Errorf("Illegal offset")
	}
	return &g.Fields[offset], nil
}

func (g *Grid) GetGroup(x, y int) (*Group, error) {
	offset := 3*x + y
	if offset > 8 {
		return nil, fmt.Errorf("Illegal offset")
	}
	return &g.Groups[offset], nil
}

func (g *Grid) Initialize() {
	for i, _ := range g.HRows {
		for j := 0; j < 9; j += 1 {
			f, err := g.GetField(j, i)
			if err != nil {
				panic(err.Error())
			}
			f.HRow = &g.HRows[i]
			g.HRows[i].Fields[j] = f
			f, err = g.GetField(i, j)
			if err != nil {
				panic(err.Error())
			}
			f.VRow = &g.VRows[i]
			g.VRows[i].Fields[j] = f
		}
		g.HRows[i].Options = getEmptyOptionsMap()
		g.VRows[i].Options = getEmptyOptionsMap()
	}
	for i, _ := range g.Groups {
		g.Groups[i].Options = getEmptyOptionsMap()
	}
	for i, _ := range g.Fields {
		yOffset := (i / 27) % 27
		xOffset := (i % 9) / 3
		gr, err := g.GetGroup(xOffset, yOffset)
		if err != nil {
			panic(err.Error())
		}
		g.Fields[i].Group = gr
		if e := gr.AppendField(&g.Fields[i]); e != nil {
			panic(e.Error())
		}
		g.Fields[i].Options = getEmptyOptionsMap()
		g.Fields[i].Grid = g
	}
	g.Options = map[int]int{}
	for k, _ := range g.Fields[0].Options {
		g.Options[k] = 9
	}
}

//use %v + this method, returns space for zero values
//returns int for set values (1 - 9)
func (f *Field) GetPrintValue() interface{} {
	var r interface{}
	if f.Value == 0 {
		r = " "
	} else {
		r = f.Value
	}
	return r
}

func (r *Row) GetValues() [9]int {
	ret := [9]int{}
	for i, f := range r.Fields {
		ret[i] = f.Value
	}
	return ret
}

func (g *Grid) String() string {
	buffer := []byte("+|===|===|===||===|===|===||===|===|===|+\n")
	for i, hr := range g.HRows {
		v := hr.GetValues()
		buffer = append(buffer, fmt.Sprintf(
			"|| %v | %v | %v || %v | %v | %v || %v | %v | %v ||\n",
			v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7], v[8], //find a way to get a []interface{} type returned, and use v... to print
		)...)
		if m := i % 3; i != 8 && m == 2 {
			buffer = append(buffer,
				"||---|---|---||---|---|---||---|---|---||\n"...,
			)
		}
	}
	buffer = append(buffer, "+|===|===|===||===|===|===||===|===|===|+"...)
	return string(buffer)
}

func (g *Grid) TestValues() *Grid {
	for k, _ := range g.Fields {
		g.Fields[k].Value = (k % 9) + 1
	}
	return g
}

func (g *Grid) SetValues(stream []byte) error {
	pos := 0
	for _, c := range stream {
		i := int(c)
		if i > 48 && i < 58 {
			if err := g.Fields[pos].SetValue(i - 48); err != nil {
				return fmt.Errorf("Field @%d: %s", pos, err.Error())
			}
		} else {
			g.Fields[pos].Value = 0
		}
		pos += 1
		if pos > 80 {
			break
		}
	}
	return nil
}

func (f *Field) SetValue(v int) error {
	if cnt, ok := f.Grid.Options[v]; !ok || cnt == 0 {
		return fmt.Errorf("%d invalid value, or all occurences used", v)
	}
	if _, ok := f.Options[v]; !ok {
		return fmt.Errorf("%d invalid value", v)
	}
	f.Value = v
	if err := f.Group.RemoveOption(v); err != nil {
		return err
	}
	if err := f.HRow.RemoveOption(v); err != nil {
		return err
	}
	if err := f.VRow.RemoveOption(v); err != nil {
		return err
	}
	for k, _ := range f.Options {
		f.Options[k] = false
	}
	f.Options[v] = true
	f.Grid.Options[v] -= 1
	return nil
}

func (g *Grid) CycleRows() bool {
	changes := false
	for v, cnt := range g.Options {
		if cnt > 0 {
			for i, _ := range g.HRows {
				g.HRows[i].TryValue(v)
				g.VRows[i].TryValue(v)
				if g.Options[v] != cnt {
					changes = true
				}
			}
		}
	}
	c2 := g.CycleGroups()
	if changes || c2 {
		return true
	}
	return false
}

func (g *Grid) CycleGroups() bool {
	changes := false
	for v, cnt := range g.Options {
		if cnt > 0 {
			for i, _ := range g.Groups {
				if set, _ := g.Groups[i].TryValue(v); set {
					changes = true
				}
			}
		}
	}
	return changes
}

func (g *Grid) IsSolved() bool {
	for _, cnt := range g.Options {
		if cnt > 0 {
			return false
		}
	}
	return true
}

func (r *Row) RemoveOption(v int) error {
	if bv, ok := r.Options[v]; !ok || !bv {
		return fmt.Errorf("%d is either an invalid value or no longer a valid option", v)
	}
	r.Options[v] = false
	for _, f := range r.Fields {
		f.Options[v] = false
	}
	return nil
}

func (gr *Group) RemoveOption(v int) error {
	if bv, ok := gr.Options[v]; !ok || !bv {
		return fmt.Errorf("%d is either invalid or no longer available", v)
	}
	for _, f := range gr.Fields {
		f.Options[v] = false
	}
	return nil
}

func (gr *Group) TryValue(v int) (bool, error) {
	bv, ok := gr.Options[v]
	if !ok {
		return false, fmt.Errorf("%d invalid value", v)
	}
	if !bv {
		return false, nil
	}
	fields := gr.getPossibleFields(v)
	if len(fields) == 1 {
		if err := fields[0].SetValue(v); err != nil {
			return false, err
		}
		return true, nil
	}
	for _, f := range fields {
		avail := f.getPossibleValues()
		if len(avail) == 1 {
			if err := f.SetValue(v); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func (gr *Group) getPossibleFields(v int) []*Field {
	ret := []*Field{}
	for _, f := range gr.Fields {
		if f.Value == 0 && f.Options[v] {
			ret = append(ret, f)
		}
	}
	return ret
}

//Return true if value was set
//returns error for invalid values
//return false if value could not be set - or was already set
func (r *Row) TryValue(v int) (bool, error) {
	bv, ok := r.Options[v]
	if !ok {
		return false, fmt.Errorf("%d invalid value", v)
	}
	if !bv {
		return false, nil
	}
	fields := r.getPossibleFields(v)
	if len(fields) == 1 {
		if err := fields[0].SetValue(v); err != nil {
			return false, err
		}
		return true, nil
	}
	for _, f := range fields {
		avail := f.getPossibleValues()
		if len(avail) == 1 {
			if err := f.SetValue(v); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func (r *Row) getPossibleFields(v int) []*Field {
	ret := []*Field{}
	for _, f := range r.Fields {
		if f.Value == 0 && f.Options[v] {
			ret = append(ret, f)
		}
	}
	return ret
}

func (f *Field) getPossibleValues() []int {
	ret := []int{}
	for k, b := range f.Options {
		if b {
			ret = append(ret, k)
		}
	}
	return ret
}

func main() {
	pathPtr := flag.String("path", "", "Path to in-file (containing sudoku puzzle)")
	strPtr := flag.String("raw", "", "String representation of sudoku to solve")

	var f []byte
	flag.Parse()

	if *pathPtr != "" {
		fs, err := ioutil.ReadFile(*pathPtr)
		if err != nil {
			panic(err.Error())
		}
		f = fs
	} else {
		if *strPtr == "" {
			fmt.Println("Either path or raw flag is required")
			flag.PrintDefaults()
			return
		}
		f = []byte(*strPtr)
	}

	g := &Grid{}
	g.Initialize()

	if err := g.SetValues(f); err != nil {
		panic(err.Error()) // input error
	} else {
		fmt.Println(g.String())
		tries := 10000
		for tries > 0 {
			if !g.CycleRows() {
				tries -= 1
			}
			if g.IsSolved() {
				break
			}
		}
		fmt.Println(g.String())
	}
}
