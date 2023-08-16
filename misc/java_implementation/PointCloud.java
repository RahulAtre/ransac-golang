import java.util.ArrayList;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.BufferedWriter;
import java.io.FileWriter;
import java.util.Scanner;
import java.util.StringTokenizer;
import java.util.Iterator;
import java.util.Random;

/**
 * @author Rahul Atre
 */

public class PointCloud {

    private ArrayList<Point3D> pointCloud; //Instance variable to store points in an ArrayList
    private String cloudName; //To store the cloud's name


    /**
     * Constructor from an XYZ file
     * 
     * @param filename -> User inputted filename to read and analyze from
     */
    public PointCloud(String filename) throws Exception {
        if(filename == null) {
            throw new FileNotFoundException("File name cannot be null");
        }

        Scanner scanner = new Scanner(new File(filename)); //Initiate a scanner for the file
        scanner.nextLine(); //Skip the first line [column names]

        pointCloud = new ArrayList<Point3D>(); //Instantiate pointCloud array

        /**
         * Parsing data from file
         */
        while(scanner.hasNext()) {
            StringTokenizer str = new StringTokenizer(scanner.nextLine());
            Point3D newPoint = new Point3D(Double.parseDouble(str.nextToken()), Double.parseDouble(str.nextToken()), Double.parseDouble(str.nextToken()));
            
            pointCloud.add(newPoint);
        }

        this.cloudName = filename.replaceAll(".xyz", ""); //We will need to know the filename of the cloud to create new files (dominant planes) corresponding to it 
    }


    /**
     * Empty Constructor that constructs an empty point cloud
     */
    public PointCloud() {
        pointCloud = new ArrayList<Point3D>(); //Instantiate pointCloud array
        cloudName = null;
    }


    /**
     * @return name of the point cloud
     */
    public String getName() {
        return this.cloudName;
    }


    /**
     * Adds a point to the point cloud
     * 
     * @param pt represents a point in 3D space
     */
    public void addPoint(Point3D pt) {
        if(pt==null) { 
            throw new NullPointerException("You cannot add a null point to this list");
        }

        pointCloud.add(pt);
    }


    /**
     * Remove a point to the point cloud
     * 
     * @param pt represents a point in 3D space that you can remove
     */
    public void removePoint(Point3D pt) {
        if(pt==null) { 
            throw new NullPointerException("You cannot remove a null point from this list");
        }

        pointCloud.remove(pt); //Since we know pt will be a part of the cloud, we don't need to check if it was removed or not
    }


    /**
     * @return a random point from the cloud
     */
    public Point3D getPoint() {
        Random random = new Random(); //Create a new Random object

        return pointCloud.get(random.nextInt(pointCloud.size())); //Get a random Point3D element from the pointCloud
    }

    /**
     * @return size of the point cloud
     */
    public int size() {
        return this.pointCloud.size();
    }


    /**
     * A save method that saves the point cloud into an XYZ file
     * 
     * @param filename -> User inputted filename (preferably a new filename) to store the point cloud
     */
    public void save(String filename) {
        try {
            BufferedWriter writer = new BufferedWriter(new FileWriter(new File(filename))); //Create a BufferWriter obj to write to the given file
            writer.write("x" + "\t" + "y" + "\t" + "z" + "\n"); //Column name

            for(int i = 0; i<this.pointCloud.size(); i++) { //For all points in point cloud -> write to the file
                writer.write(pointCloud.get(i).getX() + "\t" + pointCloud.get(i).getY() 
                            + "\t" + pointCloud.get(i).getZ() + "\n");
            } 
            writer.close(); //Close writer after all points are added

        } catch(Exception e) { //Catch exception if raised
            e.printStackTrace(); 
        }
    }   


    /**
     * @return an iterator to the points in the cloud
     */
    public Iterator<Point3D> iterator() {
        return pointCloud.iterator(); //ArrayList has implemented a built-in iterator in its class
    } 
}