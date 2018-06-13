// pictures.go

// Copyright (C) 2018  Steve Merrony

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tello

import (
	"fmt"
	"io/ioutil"
	"sort"
)

// TakePicture requests the Tello to take a JPEG snapshot.
// The process takes a little while to complete and the video may freeze
// during photography.  Sometime the Tello does not honour the request.
// The pictures are stored in the tello struct until saved by eg. SaveAllPics().
func (tello *Tello) TakePicture() (err error) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()

	tello.ctrlSeq++
	pkt := newPacket(ptSet, msgDoTakePic, tello.ctrlSeq, 0)
	tello.ctrlConn.Write(packetToBuffer(pkt))
	//log.Println("Sent take picture request")
	return nil
}

func (tello *Tello) sendFileSize() {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	tello.ctrlSeq++
	tello.ctrlConn.Write(packetToBuffer(newPacket(ptData1, msgFileSize, tello.ctrlSeq, 1)))
}

func (tello *Tello) sendFileAckPiece(done byte, fID uint16, pieceNum uint32) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	tello.ctrlSeq++
	pkt := newPacket(ptData1, msgFileData, tello.ctrlSeq, 7)
	pkt.payload[0] = done
	pkt.payload[1] = byte(fID)
	pkt.payload[2] = byte(fID >> 8)
	pkt.payload[3] = byte(pieceNum)
	pkt.payload[4] = byte(pieceNum >> 8)
	pkt.payload[5] = byte(pieceNum >> 16)
	pkt.payload[6] = byte(pieceNum >> 24)
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

func (tello *Tello) sendFileDone(fID uint16, size int) {
	tello.ctrlMu.Lock()
	defer tello.ctrlMu.Unlock()
	tello.ctrlSeq++
	pkt := newPacket(ptGet, msgFileDone, tello.ctrlSeq, 6)
	pkt.payload[0] = byte(fID)
	pkt.payload[1] = byte(fID >> 8)
	pkt.payload[2] = byte(size)
	pkt.payload[3] = byte(size >> 8)
	pkt.payload[4] = byte(size >> 16)
	pkt.payload[5] = byte(size >> 24)
	tello.ctrlConn.Write(packetToBuffer(pkt))
}

// reassembleFile reassembles a chunked file in tello.fileTemp into a contiguous byte array in tello.files
func (tello *Tello) reassembleFile() {
	var fd fileData
	tello.fdMu.Lock()
	defer tello.fdMu.Unlock()

	fd.fileType = tello.fileTemp.filetype
	fd.fileSize = tello.fileTemp.accumSize
	// we expect the pieces to be in order
	for _, p := range tello.fileTemp.pieces {
		// the chunks may not be in order, we must sort them
		if p.numChunks > 1 {
			sort.Slice(p.chunks, func(i, j int) bool {
				return int(p.chunks[i].chunkNum) < int(p.chunks[j].chunkNum)
			})
		}
		for _, c := range p.chunks {
			fd.fileBytes = append(fd.fileBytes, c.chunkData...)
		}
	}
	tello.files = append(tello.files, fd)
	tello.fileTemp = fileInternal{}
}

// NumPics returns the number of JPEG pictures we are storing in memory
func (tello *Tello) NumPics() (np int) {
	for _, f := range tello.files {
		if f.fileType == ftJPEG {
			np++
		}
	}
	return np
}

// SaveAllPics writes all JPEG pictures to disk using the given prefix
// and a generated index number. It returns the number of pictures written &/or an error.
// If there is no error, the pictures are removed from memory.
func (tello *Tello) SaveAllPics(prefix string) (np int, err error) {
	for _, f := range tello.files {
		if f.fileType == ftJPEG {
			filename := fmt.Sprintf("%s_%d.jpg", prefix, np)
			err = ioutil.WriteFile(filename, f.fileBytes, 0644)
			if err != nil {
				break
			}
			np++
		}
	}
	tello.files = nil
	return np, nil
}
