package algorithms;

public class NoEights {

    public static int smallestAmount(int low, int high){
        int lowest = countEight(low);
        for(int i=low; i<=high; i++){
            int count = countEight(i);
            if (count == 0)
                return 0;
            else if(count < lowest)
                lowest = count;
        }
        return lowest;
    }

    public static int countEight(int number){
        int count = 0;
        String eight = String.valueOf(number);
        for(char c : eight.toCharArray()){
            if (c == "8".chatAt(0)) count++;
        }
        return count;
    }

    public static void main(String[] args) {
        int low = 8808;
        int high = 8880;
        System.out.println(smallestAmount (low, high));
    }

}
