package mcio

import (
	"bufio"
	"convertfs/admutils"
	"log"
	"os"
	"regexp"
)

// HashFileFirstColumn takes a file path string as argument
// and calculates hash for the file's first column (see *.snp or *.ind file)
// File must have last empty line
func HashFileFirstColumn(path string) int32 {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var hash, thash int32
	reader := bufio.NewReader(file)
	firstColumnRegexp, _ := regexp.Compile(`\w+`)
	for {
		str, err := reader.ReadString(10)
		if err != nil {
			return hash
		}

		thash = admutils.Hashit(firstColumnRegexp.FindString(str))
		hash *= 17
		hash ^= thash
	}
}

// Calcishash takes ...
// and calculate hashes on individuals and SNPs (to compare with file values).
// (see https://github.com/DReichLab/EIG/blob/master/src/mcio.c#L2697)
func Calcishash(snppath, indivpath string) {
	// firstColumnRegexp, _ := regexp.Compile(`\w+`)
	// fmt.Println(firstColumnRegexp.FindString("          rs12124819     1        0.020242          776546 A G"))
}

/*!  \fn int calcishash(SNP **snpm, Indiv **indiv, int numsnps, int numind, int *pihash, int *pshash)
\brief Calculate hashes on individuals and SNPs (to compare with file values.)
\param snpm  Array of SNP data
\param indiv Array if individual data
\param numsnps  Number of elements in snpm
\param numind  Number of elements in indiv
\param pihash  Output parameter for indiv hash
\param pshash  Output parameter for SNP hash
Return number of SNPs plus number if individuals
*/

/*
int
calcishash (SNP **snpm, Indiv **indiv, int numsnps, int numind, int *pihash,
            int *pshash)
{
  char **arrx;
  int ihash, shash, n, num;
  int i;
  Indiv *indx;
  SNP *cupt;

  n = numind;
  ZALLOC(arrx, n, char *);

  num = 0;
  for (i = 0; i < n; i++)
    {
      indx = indiv[i];
      arrx[num] = strdup (indx->ID);
      ++num;
    }
  *pihash = hasharr (arrx, num);

  freeup (arrx, num);
  free (arrx);

  n = numsnps;
  ZALLOC(arrx, n, char *);
  num = 0;
  for (i = 0; i < n; i++)
    {
      cupt = snpm[i];
      if (cupt->isfake)
        continue;
      arrx[num] = strdup (cupt->ID);
      cupt->ngtypes = numind;
      ++num;
    }
  *pshash = hasharr (arrx, num);
  freeup (arrx, num);
  free (arrx);
  return num;

}
*/
