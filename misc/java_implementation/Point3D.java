/**
 * @author Rahul Atre
 */

public class Point3D {

    //Instance variables for the x, y, z co-ordinates of a point
    private double x, y, z;
    

    /**
     * Constructor for a single point
     * 
     * @param x is the x-cord 
     * @param y is the y-cord 
     * @param z is the z-cord
     */
    public Point3D(double x, double y, double z) {
        this.x = x;
        this.y = y;
        this.z = z;
    }


    /**
     * @return the x-value of this point
     */
    public double getX() {
        return this.x;
    }


    /**
     * @return the y-value of this point
     */
    public double getY() {
        return this.y;
    }


    /**
     * @return the z-value of this point
     */
    public double getZ() {
        return this.z;
    }


    /**
     * @return's a string representation of a Point
     */
    public String toString() {
        return "(" + this.x + ", " + this.y + ", " + this.z + ")";
    }
}
