////go:generate go run ../utils/gensnapshots/gensnapshots.go -rubi_snapshot_location=./rubi_snapshot/rubi_snapshot.expred
////go:generate go-bindata -pkg rubi_snapshot -o rubi_snapshot/rubi_resources.go -nocompress rubi_snapshot/rubi_snapshot.expred

package expreduce

func getColorDefinitions() (defs []Definition) {
	defs = append(defs, Definition{
		Name: "Black",
	})
	defs = append(defs, Definition{
		Name: "White",
	})
	defs = append(defs, Definition{
		Name: "Red",
	})
	defs = append(defs, Definition{
		Name: "Green",
	})
	defs = append(defs, Definition{
		Name: "Blue",
	})
	defs = append(defs, Definition{
		Name: "Yellow",
	})
	defs = append(defs, Definition{
		Name: "Cyan",
	})
	defs = append(defs, Definition{
		Name: "Magenta",
	})
	defs = append(defs, Definition{
		Name: "Brown",
	})
	defs = append(defs, Definition{
		Name: "Orange",
	})
	defs = append(defs, Definition{
		Name: "Pink",
	})
	defs = append(defs, Definition{
		Name: "Purple",
	})
	defs = append(defs, Definition{
		Name: "Grey",
	})
	return
}
