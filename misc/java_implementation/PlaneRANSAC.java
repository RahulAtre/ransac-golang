import java.lang.Math;
import java.util.Iterator;
import java.util.ArrayList;

/**
 * @author Rahul Atre
 */

public class PlaneRANSAC {

    private double eps; //Instance variable for epsilon distance
    private PointCloud pointCloud; //Instance variable to store the point cloud

    /**
     * Constructor that utilizes the RANSAC algorithm
     * @param pt PointCloud input from the user
     */
    public PlaneRANSAC(PointCloud pt) {
        this.pointCloud = pt;
        this.eps = 0.1; //Epsilon Value set to 0.1 for consistency purposes

        for(int threeMostDominant = 1; threeMostDominant <= 3; threeMostDominant++) { //For-loop that runs 3 times to get top three dominant planes

            String newFileExtension = "";

            switch(threeMostDominant) { //Create a new file based on which dominant plane is being found
            case 1: 
                newFileExtension = "_1.xyz"; 
                break;
            case 2: 
                newFileExtension = "_2.xyz"; 
                break;
            case 3: 
                newFileExtension = "_3.xyz"; 
                break;
            }

            String outputFileName = pt.getName() + newFileExtension; //File name
            int numOfIterations = this.getNumberOfIterations(0.99, 0.1); //Here, we have chosen the appropriate values for the confidence and probability %. For this assignment, I felt it was best to tweak the values only through the constructor for consistency purposes.

            this.run(numOfIterations, outputFileName);
        }
        pointCloud.save(pointCloud.getName() + "_p0.xyz"); //Original cloud without the plane's points

        System.out.println("Algorithm Complete"); //Output to console indicating that one algorithm run has completed
    }


    /**
     * Getter for the epsilon value
     */
    public double getEps() {
        return this.eps;
    }


    /**
     * Setter for the epsilon value
     */
    public void setEps(double eps) {
        this.eps = eps;
    }


    /**
     * @param confidence 
     * @param percentageOfPointsOnPlane
     * 
     * @return the estimated number of iterations required to obtain a certain level of confidence
     * to identify a plane made of a certain percentage of points
     */
    public int getNumberOfIterations(double confidence, double percentageOfPointsOnPlane) {
        int numberOfIterations = (int) (Math.log10(1 - confidence) / Math.log10(1 - Math.pow(percentageOfPointsOnPlane, 3))); //Equation for num of iterations

        return numberOfIterations; 
    }


    /**
     * A method that runs the RANSAC algorithm for identifying the dominant plane of the point cloud (only one plane)
     * 
     * @param numberOfIterations 
     * @param filename -> filename of most dominant plane
     */
    public void run(int numberOfIterations, String filename) {

        PointCloud bestSupport = new PointCloud(); //Best support plane 
        int iterations=0;

        while(iterations < numberOfIterations) { //while-loop to run the program until the set number of iterations

            Plane3D currentPlane = new Plane3D(pointCloud.getPoint(), pointCloud.getPoint(), pointCloud.getPoint()); //Create plane
            PointCloud currentSupport = new PointCloud(); //Variable to store the support for the current iteration
            Iterator<Point3D> cloudIterator = pointCloud.iterator();
            
            while(cloudIterator.hasNext()) { //For all points, check if it is a support point or not
                Point3D pointFromCloud = cloudIterator.next();
                double distancePointToPlane = currentPlane.getDistance(pointFromCloud);

                if(distancePointToPlane < eps) {
                    currentSupport.addPoint(pointFromCloud);
                }
            }

            if(currentSupport.size() > bestSupport.size()) { //If the # of support points on this plane is bigger than the best support, current will replace best support
                bestSupport = currentSupport;
            }
            iterations++;
        }

        bestSupport.save(filename); //Creation of dominant plane file  

        Iterator<Point3D> bestSupportIterator = bestSupport.iterator(); 
        while(bestSupportIterator.hasNext()) {
            pointCloud.removePoint(bestSupportIterator.next()); //Keep removing support points from point cloud until the array is empty
        }
    }


    public static void main(String[] args) {
        try {

            PointCloud pointCloud1 = new PointCloud("PointCloud1.xyz"); //Obtain pointcloud from file
            PlaneRANSAC ransacAlgorithm = new PlaneRANSAC(pointCloud1); //Run RANSAC algorithm on PointCloud1

            PointCloud pointCloud2 = new PointCloud("PointCloud2.xyz"); //Obtain pointcloud from file
            PlaneRANSAC ransacAlgorithm2 = new PlaneRANSAC(pointCloud2); //Run RANSAC algorithm on PointCloud2

            PointCloud pointCloud3 = new PointCloud("PointCloud3.xyz"); //Obtain pointcloud from file
            PlaneRANSAC ransacAlgorithm3 = new PlaneRANSAC(pointCloud3); //Run RANSAC algorithm on PointCloud3
        
        } catch(Exception e) {
            System.out.println("Something went wrong while trying to open/read the file");
        }
    }
}