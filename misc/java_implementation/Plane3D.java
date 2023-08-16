import java.lang.Math;
import java.util.ArrayList;

/**
 * @author Rahul Atre
 */

public class Plane3D {

    //Instance variables for a, b, c, d constant values of a plane 
    private double a, b, c, d;

    //Instance variables for plane points if a plane is constructed from 3 points
    private Point3D[] planePoints;

    /**
     * Constructor from 3 points
     * 
     * @param p1 is the 1st Point -> Use 3 points to represent a plane
     * @param p2 is the 2nd Point
     * @param p3 is the 3rd Point
     */
    public Plane3D(Point3D p1, Point3D p2, Point3D p3) {
        /**
         * Algorithm for finding eq. of a plane given 3 points
         * 1. Create two vectors by subtracting x,y,z coordinates from any two points
         * 2. Take the cross product of the two vectors -> Gives normal vector w/ a, b, c
         * 3. Plug any point into the equation of a plane to find d [a(x) + b(y) + c(z) = d]
         */
        double x1 = p1.getX() - p2.getX();
        double y1 = p1.getY() - p2.getY();
        double z1 = p1.getZ() - p2.getZ();

        double x2 = p1.getX() - p3.getX();
        double y2 = p1.getY() - p3.getY();
        double z2 = p1.getZ() - p3.getZ();

        this.a = y1 * z2 - y2 * z1;
        this.b = x2 * z1 - x1 * z2;
        this.c = x1 * y2 - x2 * y1;
        this.d = -(a*p1.getX() + b*p1.getY() + c*p1.getZ());

        this.planePoints = new Point3D[3];
        this.planePoints[0] = p1;
        this.planePoints[1] = p2;
        this.planePoints[2] = p3;
    }
    

    /**
     * Constructor from plane parameters
     * 
     * @param a is the 1st Constant -> Use 4 constant variables to represent the equation of a plane [ax + by + cz + d = 0]
     * @param b is the 2nd Constant 
     * @param c is the 3rd Constant 
     * @param d is the 4th Constant
     */
    public Plane3D(double a, double b, double c, double d) {
        this.a = a;
        this.b = b;
        this.c = c;
        this.d = d;

        planePoints = new Point3D[3];
    }


    /**
     * @param index specifies which point the user would like from the plane
     * @return The specificed point from the plane
     */
    public Point3D getPoint(int index) {
        if(index > 2 || index < 0) {
            throw new IndexOutOfBoundsException("Invalid index, please enter a valid array slot");
        }

        return planePoints[index];
    }
    

    /**
     * @param pt represents a point in 3D space
     * @return the distance from pt (point) to the plane
     */
    public double getDistance(Point3D pt) {
        //Formula for calculating the distance from a point to a plane
        double numerator = Math.abs(this.a*pt.getX() + this.b*pt.getY() + this.c*pt.getZ() + this.d);
        double denominator = Math.sqrt(Math.pow(this.a, 2) + Math.pow(this.b, 2) + Math.pow(this.c, 2));

        return numerator/denominator;
    }
}