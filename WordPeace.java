package algorithms;

import java.util.Arrays;

public class WordPeace {

    public static long numGroups(int k, int[] countries) {
        long num = 0;
        while (true) {
            Arrays.sort(countries);
            int[] reversed = new int[countries.length];
            for (int i = countries.length - 1, j = 0; i >= 0; i--, j++) {
                reversed[j] = countries[i];
            }
            countries = reversed;
            boolean canForm = false;
            int min = min(countries, k);
            int last = -1;
            for (int i = 0; i < k; i++) {
                last = getNextCountry(last, countries, min);
                if (last == -1) {
                    canForm = false;
                    break;
                } else {
                    canForm = true;
                }
            }
            if (canForm) {
                num += min;
            } else {
                break;
            }
        }
        return num;
    }

    public static int getNextCountry(int last, int[] countries, int min) {
        for (int i = last + 1; i < countries.length; i++) {
            if (countries[i] >= min) {
                countries[i] -= min;
                return i;
            }
        }
        return -1;
    }

    public static int min(int[] array, int k) {
        int max = array[array.length - k];
        int min = 0;
        for (int i = array.length - 1; i >= 0; i--) {
            if (array[i] == 0) {
                continue;
            } else {
                min = array[i];
                break;
            }
        }
        if (min <= 20) {
            return 1;
        } else {
            return min = min == max ? min / 2 : min;
        }
    }

    public static void main(String[] args) {
        int k = 5;
        int[] countries = {1, 2, 3, 4, 5, 6};
        System.out.println(numGroups(k, countries));
    }
}
