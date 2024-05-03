/*
165. Compare Version Numbers
Given two version numbers, version1 and version2, compare them.

Version numbers consist of one or more revisions joined by a dot '.'. Each revision consists of 
digits and may contain leading zeros. Every revision contains at least one character. Revisions 
are 0-indexed from left to right, with the leftmost revision being revision 0, the next revision 
being revision 1, and so on. For example 2.5.33 and 0.1 are valid version numbers.

To compare version numbers, compare their revisions in left-to-right order. Revisions are 
compared using their integer value ignoring any leading zeros. This means that revisions 1 
and 001 are considered equal. If a version number does not specify a revision at an index, 
then treat the revision as 0. For example, version 1.0 is less than version 1.1 because their 
revision 0s are the same, but their revision 1s are 0 and 1 respectively, and 0 < 1.

Return the following:

If version1 < version2, return -1.
If version1 > version2, return 1.
Otherwise, return 0.
 */

public class CompareVersion {
    public static int compareVersion(String version1, String version2) {
        String[] v1 = version1.split("\\.");
        String[] v2 = version2.split("\\.");

        int maxLength = Math.max(v1.length, v2.length);
        int res = 0;
        for (int i = 0; i < maxLength; i++) {
            int num1 = i < v1.length ? Integer.parseInt(v1[i]) : 0;
            int num2 = i < v2.length ? Integer.parseInt(v2[i]) : 0;

            if (num1 < num2) {
                res = -1;
                break;
            } else if (num1 > num2) {
                res = 1;
                break;
            }
        }
        return res;
    }

    public static void main(String[] args) {
        System.out.println(compareVersion("1.05", "1.1"));
    }
}
