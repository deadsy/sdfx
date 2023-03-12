package mesh

import (
	"fmt"
	"os"
	"time"

	"github.com/deadsy/sdfx/render/finiteelements/buffer"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Inp writes different types of finite elements as ABAQUS or CalculiX `inp` file.
type Inp struct {
	// Finite elements mesh.
	Mesh FE
	// Output `inp` file path.
	Path string
	// For writing nodes to a separate file.
	PathNodes string
	// For writing elements to a separate file.
	PathEls string
	// For writing boundary conditions to a separate file.
	PathBou string
	// Output `inp` file would include start layer.
	LayerStart int
	// Output `inp` file would exclude end layer.
	LayerEnd int
	// Layers fixed to the 3D print floor i.e. bottom layers. The boundary conditions.
	LayersFixed []int
	// To write only required nodes to `inp` file.
	TempVBuff *buffer.VB
	// Mechanical properties of 3D print resin.
	MassDensity  float32
	YoungModulus float32
	PoissonRatio float32
}

// NewInp sets up a new writer.
func NewInp(
	m FE,
	path string,
	layerStart, layerEnd int,
	layersFixed []int,
	massDensity float32, youngModulus float32, poissonRatio float32,
) *Inp {
	return &Inp{
		Mesh:         m,
		Path:         path,
		PathNodes:    path + ".nodes",
		PathEls:      path + ".elements",
		PathBou:      path + ".boundary",
		LayerStart:   layerStart,
		LayerEnd:     layerEnd,
		LayersFixed:  layersFixed,
		TempVBuff:    buffer.NewVB(),
		MassDensity:  massDensity,
		YoungModulus: youngModulus,
		PoissonRatio: poissonRatio,
	}
}

// Write starts writing to `inp` file.
func (inp *Inp) Write() error {
	f, err := os.Create(inp.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = inp.writeHeader(f)
	if err != nil {
		return err
	}

	// Write nodes.

	_, err = f.WriteString("*NODE\n")
	if err != nil {
		return err
	}

	// Include a separate file to avoid cluttering the `inp` file.
	_, err = f.WriteString(fmt.Sprintf("*INCLUDE,INPUT=%s\n", inp.PathNodes))
	if err != nil {
		return err
	}

	// Write to a separate file to avoid cluttering the `inp` file.
	fNodes, err := os.Create(inp.PathNodes)
	if err != nil {
		return err
	}
	defer fNodes.Close()

	// Temp buffer is just to avoid writing repeated nodes into the `inpt` file.
	defer inp.TempVBuff.DestroyHashTable()

	err = inp.writeNodes(fNodes)
	if err != nil {
		return err
	}

	// Write elements.

	ElementType := ""
	if inp.Mesh.Npe() == 4 {
		ElementType = "C3D4"
	} else if inp.Mesh.Npe() == 8 {
		ElementType = "C3D8"
	} else if inp.Mesh.Npe() == 20 {
		ElementType = "C3D20R"
	}

	_, err = f.WriteString(fmt.Sprintf("*ELEMENT, TYPE=%s, ELSET=Eall\n", ElementType))
	if err != nil {
		return err
	}

	// Include a separate file to avoid cluttering the `inp` file.
	_, err = f.WriteString(fmt.Sprintf("*INCLUDE,INPUT=%s\n", inp.PathEls))
	if err != nil {
		return err
	}

	// Write to a separate file to avoid cluttering the `inp` file.
	fEls, err := os.Create(inp.PathEls)
	if err != nil {
		return err
	}
	defer fEls.Close()

	err = inp.writeElements(fEls)
	if err != nil {
		return err
	}

	// Fix the degrees of freedom one through three for all nodes on specific layers.

	_, err = f.WriteString("*BOUNDARY\n")
	if err != nil {
		return err
	}

	// Include a separate file to avoid cluttering the `inp` file.
	_, err = f.WriteString(fmt.Sprintf("*INCLUDE,INPUT=%s\n", inp.PathBou))
	if err != nil {
		return err
	}

	// Write to a separate file to avoid cluttering the `inp` file.
	fBou, err := os.Create(inp.PathBou)
	if err != nil {
		return err
	}
	defer fBou.Close()

	err = inp.writeBoundary(fBou)
	if err != nil {
		return err
	}

	return inp.writeFooter(f)
}

func (inp *Inp) writeHeader(f *os.File) error {
	if 0 <= inp.LayerStart && inp.LayerStart < inp.LayerEnd && inp.LayerEnd <= inp.Mesh.layerCount() {
		// Good.
	} else {
		return fmt.Errorf("start or end layer is beyond range")
	}

	_, err := f.WriteString("**\n** Structure: finite elements of a 3D model.\n** Generated by: https://github.com/deadsy/sdfx\n**\n")
	if err != nil {
		return err
	}

	_, err = f.WriteString("*HEADING\nModel: 3D model Date: " + time.Now().UTC().Format("2006-Jan-02 MST") + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (inp *Inp) writeNodes(f *os.File) error {
	// Declare vars outside loop for efficiency.
	var err error
	nodes := make([]v3.Vec, 0, inp.Mesh.Npe())
	ids := make([]uint32, inp.Mesh.Npe())
	for l := inp.LayerStart; l < inp.LayerEnd; l++ {
		for i := 0; i < inp.Mesh.feCountOnLayer(l); i++ {
			// Get the node IDs.
			nodes = inp.Mesh.feVertices(l, i)
			for n := 0; n < inp.Mesh.Npe(); n++ {
				ids[n] = inp.TempVBuff.Id(nodes[n])
			}

			// Write the node IDs.
			for n := 0; n < inp.Mesh.Npe(); n++ {
				// ID starts from one not zero.
				_, err = f.WriteString(fmt.Sprintf("%d,%f,%f,%f\n", ids[n]+1, float32(nodes[n].X), float32(nodes[n].Y), float32(nodes[n].Z)))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (inp *Inp) writeElements(f *os.File) error {
	// Declare vars outside loop for efficiency.
	var err error
	nodes := make([]v3.Vec, 0, inp.Mesh.Npe())
	ids := make([]uint32, inp.Mesh.Npe())
	var eleID uint32
	for l := inp.LayerStart; l < inp.LayerEnd; l++ {
		for i := 0; i < inp.Mesh.feCountOnLayer(l); i++ {
			nodes = inp.Mesh.feVertices(l, i)
			for n := 0; n < inp.Mesh.Npe(); n++ {
				ids[n] = inp.TempVBuff.Id(nodes[n])
			}

			// ID starts from one not zero.

			if inp.Mesh.Npe() == 4 {
				_, err = f.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d\n", eleID+1, ids[0]+1, ids[1]+1, ids[2]+1, ids[3]+1))
			} else if inp.Mesh.Npe() == 8 {
				_, err = f.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d\n", eleID+1, ids[0]+1, ids[1]+1, ids[2]+1, ids[3]+1, ids[4]+1, ids[5]+1, ids[6]+1, ids[7]+1))
			} else if inp.Mesh.Npe() == 20 {
				// There should not be more than 16 entries in a line;
				// That's why there is new line in the middle.
				// Refer to CalculiX solver documentation:
				// http://www.dhondt.de/ccx_2.20.pdf
				_, err = f.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,\n%d,%d,%d,%d,%d\n", eleID+1, ids[0]+1, ids[1]+1, ids[2]+1, ids[3]+1, ids[4]+1, ids[5]+1, ids[6]+1, ids[7]+1, ids[8]+1, ids[9]+1, ids[10]+1, ids[11]+1, ids[12]+1, ids[13]+1, ids[14]+1, ids[15]+1, ids[16]+1, ids[17]+1, ids[18]+1, ids[19]+1))
			}

			if err != nil {
				return err
			}
			eleID++
		}
	}

	return nil
}

func (inp *Inp) writeBoundary(f *os.File) error {
	// Declare vars outside loop for efficiency.
	var err error
	nodes := make([]v3.Vec, 0, inp.Mesh.Npe())
	ids := make([]uint32, inp.Mesh.Npe())
	for l := range inp.LayersFixed {
		for i := 0; i < inp.Mesh.feCountOnLayer(l); i++ {
			nodes = inp.Mesh.feVertices(l, i)
			for n := 0; n < inp.Mesh.Npe(); n++ {
				ids[n] = inp.TempVBuff.Id(nodes[n])
			}

			// Write the node IDs.
			for n := 0; n < inp.Mesh.Npe(); n++ {
				// ID starts from one not zero.
				_, err = f.WriteString(fmt.Sprintf("%d,1,3\n", ids[n]+1))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (inp *Inp) writeFooter(f *os.File) error {

	// Define material.
	// Units of measurement are mm,N,s,K.
	// Refer to:
	// https://engineering.stackexchange.com/q/54454/15178
	// Refer to:
	// Units chapter of CalculiX solver documentation:
	// http://www.dhondt.de/ccx_2.20.pdf

	_, err := f.WriteString("*MATERIAL, name=resin\n")
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("*ELASTIC,TYPE=ISO\n%e,%e,0\n", inp.YoungModulus, inp.PoissonRatio))
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("*DENSITY\n%e\n", inp.MassDensity))
	if err != nil {
		return err
	}

	// Assign material to all elements
	_, err = f.WriteString("*SOLID SECTION,MATERIAL=resin,ELSET=Eall\n")
	if err != nil {
		return err
	}

	// Write analysis

	_, err = f.WriteString("*STEP\n*STATIC\n")
	if err != nil {
		return err
	}

	// Write distributed loads.

	_, err = f.WriteString("*DLOAD\n")
	if err != nil {
		return err
	}

	// Assign gravity loading in the "positive" z-direction with magnitude 9810 to all elements.
	//
	// SLA 3D printing is done upside-down. 3D model is hanging from the print floor.
	// That's why gravity is in "positive" z-direction.
	// Here ”gravity” really stands for any acceleration vector.
	//
	// Refer to CalculiX solver documentation:
	// http://www.dhondt.de/ccx_2.20.pdf
	_, err = f.WriteString("Eall,GRAV,9810.,0.,0.,1.\n")
	if err != nil {
		return err
	}

	// Pick element results.

	_, err = f.WriteString("*EL FILE\n")
	if err != nil {
		return err
	}

	_, err = f.WriteString("S\n")
	if err != nil {
		return err
	}

	// Pick node results.

	_, err = f.WriteString("*NODE FILE\n")
	if err != nil {
		return err
	}

	_, err = f.WriteString("U\n")
	if err != nil {
		return err
	}

	// Conclude.

	_, err = f.WriteString("*END STEP\n")
	if err != nil {
		return err
	}

	return nil
}
