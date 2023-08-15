/**
* Author: Rahul Atre
* RANSAC Algorithm (implemented using a concurrent pipeline): Detects the dominant planes in a cloud of 3D points. The dominant plane of a set of 3D points is the plane that contains the largest 
* number of points. A 3D point is contained in a plane if it is located at a distance less than eps (ε) from that plane. 
*/
package main

import (
	"fmt"
	"strings"
	"strconv"	
	"os"
	"bufio"
	"math"
	"math/rand"
	"time"
	"sync"
	"runtime"
)

type Point3D struct { //Struct for Point3D (x,y,z)
	X float64
	Y float64
	Z float64
}

type Plane3D struct { //Struct for Plane3D (a,b,c,d)
	A float64
	B float64
	C float64
	D float64
}

type Plane3DwSupport struct { //Struct for Plane3DwSupport, containing the Plane & number of support points
	Plane3D
	SupportSize int
}

// reads an XYZ file and returns a slice of Point3D
func ReadXYZ(filename string) []Point3D {

	file, err := os.Open(filename) //Obtain the file from directory
	if err != nil { //Error handling
		fmt.Println("Could not open the file, here is the error message:",err.Error())
		return []Point3D{} //Null value
	}

	defer file.Close() //Close file once method is executed

	scanner := bufio.NewScanner(file) //Create a scanner to read file
	scanner.Scan(); scanner.Text() //Skip the first line [column names]

	var pointCloud []Point3D = make([]Point3D, 0) //Declare a pointCloud with initial capacity of 0

	for scanner.Scan() {
		lineInformation := strings.Split(scanner.Text(), "\t") //Grab line from file

		xCord, err := strconv.ParseFloat(lineInformation[0], 64) //Obtain (x,y,z) co-ordinates
		yCord, err := strconv.ParseFloat(lineInformation[1], 64)
		zCord, err := strconv.ParseFloat(lineInformation[2], 64)

		if err != nil { //Error handling
			fmt.Println("Could not convert string value to float, here is the error message:",err.Error())
			return []Point3D{} //Null value
		}
		pointCloud = append(pointCloud, Point3D{X: xCord, Y: yCord, Z: zCord}) //Append point to pointcloud
	}

	return pointCloud
}


// saves a slice of Point3D into an XYZ file
func SaveXYZ(filename string, points []Point3D) {
	
	file, err := os.Create(filename) //Create a new file based on provided name for file
	if err != nil { //Error handling
		fmt.Println("Could not create the file, here is the error message:",err.Error())
		return 
	}

	defer file.Close() //Close file after information has been written inside

	file.WriteString("x" + "\t" + "y" + "\t" + "z" + "\n") //Initial column header
	for _, point := range points { //For all points -> add information in the form of 'string' and write data to file

		lineData := strconv.FormatFloat(point.X, 'f', -1, 64) + "\t" + strconv.FormatFloat(point.Y, 'f', -1, 64) + "\t" + strconv.FormatFloat(point.Z, 'f', -1, 64) + "\n"
		file.WriteString(lineData)
	}

}


// computes the distance between point p1 and plane
func GetDistance(p1 *Point3D, plane Plane3D) float64 {
	//Formula for calculating the distance from a point to a plane
    var numerator float64 = math.Abs(plane.A*p1.X + plane.B*p1.Y + plane.C*p1.Z + plane.D);
    var denominator float64 = math.Sqrt(math.Pow(plane.A, 2) + math.Pow(plane.B, 2) + math.Pow(plane.C, 2));

    return numerator/denominator;
}

// computes the plane defined by a set of 3 points
func GetPlane(points [3]Point3D) Plane3D {
	/**
	* Algorithm for finding eq. of a plane given 3 points
	* 1. Create two vectors by subtracting x,y,z coordinates from any two points
	* 2. Take the cross product of the two vectors -> Gives normal vector w/ a, b, c
	* 3. Plug any point into the equation of a plane to find d [a(x) + b(y) + c(z) = d]
	*/

	var x1 float64 = points[0].X - points[1].X;
	var y1 float64 = points[0].Y - points[1].Y;
	var z1 float64 = points[0].Z - points[1].Z;

	var x2 float64 = points[0].X - points[2].X;
	var y2 float64 = points[0].Y - points[2].Y;
	var z2 float64 = points[0].Z - points[2].Z;

	var a float64 = y1 * z2 - y2 * z1
	var b float64 = x2 * z1 - x1 * z2
	var c float64 = x1 * y2 - x2 * y1
	var d float64 = -(a*points[0].X + b*points[0].Y + c*points[0].Z)

	return Plane3D{A: a, B: b, C: c, D: d}
}

// computes the number of required RANSAC iterations
func GetNumberOfIterations(confidence float64, percentageOfPointsOnPlane float64) int {
	//Equation for num of iterations
	var numberOfIterations int = int(math.Log10(1 - confidence) / math.Log10(1 - math.Pow(percentageOfPointsOnPlane, 3))); 

	return numberOfIterations
}

// computes the support of a plane in a set of points
func GetSupport(plane Plane3D, points []Point3D, eps float64) Plane3DwSupport {

	var numOfSupportPoints int 
	var distancePointToPlane float64
	
	for _, point := range points { //For all points in point cloud -> check if the distance between the point and plane is less than eps (so that point inside cloud)
		distancePointToPlane = GetDistance(&point, plane)

		if(distancePointToPlane < eps) {
			numOfSupportPoints++
		}
	}
	var planeSupport Plane3DwSupport = Plane3DwSupport{plane, numOfSupportPoints} //Append Plane + support size

	return planeSupport
}

// extracts the points that supports the given plane 
// and returns them as a slice of points
func GetSupportingPoints(plane Plane3D, points []Point3D, eps float64) []Point3D {

	var supportPoints []Point3D = make([]Point3D, 0) //Initialize support point array
	var distancePointToPlane float64

	for _, point := range points { //For all points in our data set
		distancePointToPlane = GetDistance(&point, plane) //Distance of point-plane

		if(distancePointToPlane < eps) { //If the distance is less than eps -> the point is in the plane
			supportPoints = append(supportPoints, point) //Add to slice
		}
	}

	return supportPoints
}

// creates a new slice of points in which all points
// belonging to the plane have been removed
func RemovePlane(plane Plane3D, points []Point3D, eps float64) []Point3D {
	
	var pointsNotInPlane []Point3D = make([]Point3D, 0) //We will store all points NOT in the plane (Instead of removing from the plane)
	var distancePointToPlane float64

	for _, point := range points { //For all points in our data set
		distancePointToPlane = GetDistance(&point, plane) //Distance of point-plane

		if(distancePointToPlane > eps) { //If the distance is greater than eps -> the point is not in the plane
			pointsNotInPlane = append(pointsNotInPlane, point) //Add to slice
		}
	}
	
	return pointsNotInPlane
}

// Randomly selects a point from point cloud, channel transmits instances of Point3D
func randomPointGenerator(points []Point3D) <-chan Point3D {
    out := make(chan Point3D) //Channel for sending a continous stream of points
    go func() {
        for {
        	var randomIndex int = rand.Intn(len(points)) //Generate random index for pointcloud
   			out <- points[randomIndex]; //Get a random Point3D element from the pointCloud
        }
        close(out)
    }()
    return out
}


// Reads Point3D instances from input channel to accumulate 3 points. Output channel transmits arrays of Point3D (composed of three points)
func tripletOfPointGenerator(pointReceiver <-chan Point3D) <-chan [3]Point3D {
    out := make(chan [3]Point3D) //Channel to store triplet points
    go func() { 
        for {
            var tripletPoint [3]Point3D = [3]Point3D{<- pointReceiver, <- pointReceiver, <- pointReceiver} //Receive 3 points from pointReceiver channel
            out <- tripletPoint //Send triplet to triplet Channel 
        }
        close(out)
    }()
    return out
}  

// Reads arrays of Point3D and resend them. Automatically stops the pipeline after having received N arrays
func takeN(tripletReceiver <-chan[3]Point3D, nArrays int) <-chan [3]Point3D {
	out := make(chan [3]Point3D) //Channel to send triplet points
    go func() {
    	defer close(out)
        for i:=0; i < nArrays; i++ { //For capacity n, transmit triplet points from generator to takeN channel
        	triplet := <- tripletReceiver
            out <- triplet
        }
    }()
    return out
}

// Reads arrays of three Point3D, computes plane defined by these points. Output channel transmits Plane3D instances
func planeEstimator(tripletReceiver <- chan[3]Point3D) <- chan Plane3D {
	out := make(chan Plane3D) //Channel to transmit plane instances
    go func() {
        for triplet := range tripletReceiver { //For triplet points received, send plane to plane channel
            out <- GetPlane(triplet)
        }
        close(out)
    }()
    return out
}

// Counts number of points in point cloud that supports the received 3D plane. Output channel transmits the plane parameters and the number of supporting 
// points in a Point3DwSupport instance
func supportPointFinder(plane Plane3D, pointCloud []Point3D, eps float64) chan Plane3DwSupport {
	out := make(chan Plane3DwSupport, 1) //Channel to trasmit Plane3DwSupport instances
    go func() {
        out <- GetSupport(plane, pointCloud, eps) //Sent from GetSupport method
        close(out)
    }()
    return out
}

// Multiplexes results received from multiple channels into one output channel
func fanIn(supportChannels []chan Plane3DwSupport) <-chan Plane3DwSupport { //Source: GO Concurrency Pipeline documentation https://go.dev/blog/pipelines
    var wg sync.WaitGroup
    out := make(chan Plane3DwSupport)

    // Start an output goroutine for each input channel in supportChannels.
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan Plane3DwSupport) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(supportChannels))
    for _, c := range supportChannels {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

// Receives Plane3DwSupport instances & keeps in memory the plane w/ best support received so far. This component does not output values,
// simply maintains the provided *Plane3DwSupport variable
func dominantPlaneIdentifier(plane3DSupportReceiver <- chan Plane3DwSupport, bestSupport *Plane3DwSupport) {
	for plane3DwSupport := range plane3DSupportReceiver { //For all planes in channel -> if the current plane has a higher support than the best support plane, it becomes the new best support plane
		if plane3DwSupport.SupportSize > bestSupport.SupportSize {
			*bestSupport = plane3DwSupport
		}
	}
}


func pipeline(confidenceNum float64, percentageNum float64, eps float64, pointCloud []Point3D, numOfIterations int) ([]Point3D, []Point3D) {
	
	var bestSupport Plane3DwSupport	= Plane3DwSupport{Plane3D{0,0,0,0}, 0} //Instantiate best support instance

	//Create pipeline for the given pointCloud
	randomPointChannel := randomPointGenerator(pointCloud) 
	tripletChannel := tripletOfPointGenerator(randomPointChannel)
	takeNChannel := takeN(tripletChannel, numOfIterations)
	planeChannel := planeEstimator(takeNChannel)

	var numberOfSupportFinderChannels int = numOfIterations
	supportArrayChan := make([]chan Plane3DwSupport, numOfIterations) //Create an array of support channels

	var wg sync.WaitGroup //WaitGroup to synchronize go routines

	for i := 0; i < numberOfSupportFinderChannels; i++ { //For all Plane3DwSupport channels, create a new support point finder channel
	    wg.Add(1) //increment for each goroutine

	    go func(i int) { //Go-routine for concurrent behavior
	        defer wg.Done()

	        supportArrayChan[i] = supportPointFinder(<-planeChannel, pointCloud, eps) //Distribute work for supportPointFinder across multiple channels
	    }(i)
	}
	wg.Wait() 

	
	bigSupportChannel := fanIn(supportArrayChan) //Fan-in process to merge all channels into one big support channel
	dominantPlaneIdentifier(bigSupportChannel, &bestSupport) //Find the dominant plane from the given unified support channel

	var dominantPlaneSupportPoints []Point3D = GetSupportingPoints(bestSupport.Plane3D, pointCloud, eps) //Return dominant plane support points
	pointCloud = RemovePlane(bestSupport.Plane3D, pointCloud, eps) //Return new pointCloud with points not from best support plane 

	return pointCloud, dominantPlaneSupportPoints
}


func main() {
	if len(os.Args) <= 4 {
		panic("To run the program, you must provide the required arguments. Please try again. Sample Argument: go run planeRANSAC.go PointCloud1.xyz 0.99 0.2 0.5")
	}

	runtime.GOMAXPROCS(8) //Setting the max number of threads that can be used 
	rand.Seed(time.Now().UnixNano()) //Change the seed each time the program is run to generate truly random numbers

	startTime := time.Now() //Measure the amount of time it took to run the program

	//Grab data from the main terminal | Parse from String to Float64 for numerical values
	var filename string = os.Args[1]
	confidenceNum, _ := strconv.ParseFloat(os.Args[2], 64)
	percentageNum, _ := strconv.ParseFloat(os.Args[3], 64)
	eps, _ := strconv.ParseFloat(os.Args[4], 64)

	var numOfIterations int = GetNumberOfIterations(confidenceNum, percentageNum) //Obtain number of iterations for RANSAC based on parameters given
	pointCloud := ReadXYZ(filename) //Read pointCloud data from file 

	var outputFileName string = strings.Replace(filename, ".xyz", "", -1) //File name without ".xyz" (used later)
	var dominantPlaneSupportPoints []Point3D //To save the support points of the 3 dominant planes in an xyz file later


	for threeMostDominant := 1; threeMostDominant <= 3; threeMostDominant++ { //For-loop that runs 3 times to get top three dominant planes

    	var newFileExtension string = "";

        switch(threeMostDominant) { //Create a new file based on which dominant plane is being found
        	case 1: 
            	newFileExtension = "_p1.xyz" 
            	break
            case 2: 
                newFileExtension = "_p2.xyz"
                break
            case 3: 
                newFileExtension = "_p3.xyz"
                break
            }

		pointCloud, dominantPlaneSupportPoints = pipeline(confidenceNum, percentageNum, eps, pointCloud, numOfIterations)
		SaveXYZ(outputFileName + newFileExtension, dominantPlaneSupportPoints) //Save dominant plane support points
    }
      
        SaveXYZ(outputFileName + "_p0.xyz", pointCloud); //Original cloud without the plane's points

        endTime := time.Since(startTime) //Calculate endtime after program finished
        fmt.Println("Algorithm Complete. Program runtime:", endTime); //Output to console indicating that one algorithm run has completed

}
