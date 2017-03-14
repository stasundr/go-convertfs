package main

import (
	"convertfs/mcio"
	"flag"
	"fmt"
	"math"
)

func main() {
	var parFile string
	flag.StringVar(&parFile, "p", "", "par file")
	flag.Parse()
	fmt.Println(parFile)

	hashOk := mcio.Calcishash(
		"/Users/me/Desktop/v19/v19.0_HO.pruned.geno",
		"/Users/me/Desktop/v19/v19.0_HO.pruned.ind",
		"/Users/me/Desktop/v19/v19.0_HO.pruned.snp")

	indNum := 435
	chunkSize := int(math.Ceil(float64(indNum) / 4))
	fmt.Println(hashOk, chunkSize)

	// https://github.com/stasundr/co-huge-converter/blob/master/modules/utils.js
	// const indNum = meta.split(/\s+/)[1];
	//       const chunkSize = Math.ceil(indNum/4);

	//       // 11000000 - 192
	//       // 00110000 - 48
	//       // 00001100 - 12
	//       // 00000011 - 3
	//       // const mask = [192, 48, 12, 3];

	//       let stream = fs.createReadStream(path, { encoding: null, start: chunkSize });
	//       stream.on('error', reject);
	//       stream.on('readable', () => {
	//         let chunk;
	//         while ((chunk = stream.read(chunkSize)) !== null) {
	//           //let snp = (chunk.readUInt8(Math.ceil(index/4)) << ((index % 4) * 2)) & 192;
	//           let snp = (chunk[Math.floor(index/4)] >> ((3 - index % 4) * 2)) & 3;

	//           byte = byte | (snp << (6 - bitpair * 2));
	//           bitpair++;

	//           if (bitpair == 4) {
	//             geno.push(byte);
	//             byte = 0;
	//             bitpair = 0;
	//           }
	//         }
	//       });
	//       stream.on('close', () => resolve(geno));
}
