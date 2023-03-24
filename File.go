package et

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

/**
文件操作工具类
 */

const FILE_PERM_DEFAULT = 0755
const FILE_TEMP_BUF_DEFAULT = 4096
const FILE_COPY_ACTION_SKIP_EXISTED = "skip_existed"
const FILE_COPY_ACTION_OVERWRITE_EXISTED = "ovewrite_existed"
var FILE_BOM_UTF8 []byte = []byte{0xef,0xbb,0xbf}
var FILE_BOM_UNICODE []byte = []byte{0xff,0xfe}


/**
文件大小
 */
func FileSize(strFilePath string) int64 {
	pFile,err := os.Open(strFilePath)
	if err != nil{
		//DbgErr().Println("ReaderFileSize failed to access:" + strFilePath)
		return -1
	}

	statInfo,err := pFile.Stat()
	if err != nil{
		//DbgErr().Println("ReaderFileSize failed to stat:" + strFilePath)
		return -1
	}
	iFileSize := statInfo.Size()

	return iFileSize
}

/**
把文件内容读取为字符串
 */
func FileReadStr(strFilePath string) string{
	strResult := ""

	if strFilePath == ""{
		return ""
	}

	bufResult,err := ioutil.ReadFile(strFilePath)
	if err != nil{
		return ""
	}
	strResult = string(bufResult)

	return strResult
}

/**
把字符串写入文件
 */
func FileWriteStr(strFilePath string, strContent string) error{
	if strFilePath == ""{
		return errors.New("file path empty")
	}

	err := fileEnsureParentDir(strFilePath)
	if err != nil{
		return err
	}
	err = ioutil.WriteFile(strFilePath,[]byte(strContent),FILE_PERM_DEFAULT)
	if err != nil{
		return err
	}

	return nil
}

/**
把字节流写入文件
 */
func FileWriteBuf(strFilePath string, bufContent []byte) error{
	if strFilePath == ""{
		return errors.New("file path empty")
	}

	err := fileEnsureParentDir(strFilePath)
	if err != nil{
		return err
	}
	err = ioutil.WriteFile(strFilePath,bufContent,FILE_PERM_DEFAULT)
	if err != nil{
		return err
	}

	return nil
}

/**
确保上层目录已经被创建
 */
func fileEnsureParentDir(strPath string) error{
	if strPath == ""{
		return nil
	}

	iLastSlash := strings.LastIndex(strPath,"/")
	if iLastSlash <= 0{
		return nil
	}

	strDir := strPath[0:iLastSlash]
	err := os.MkdirAll(strDir,0644)
	return err
}

/**
在文件后追加字符串
 */
func FileAppendStr(strFilePath string, strContent string) error{
	return FileAppendBuf(strFilePath,[]byte(strContent))
}

/**
在文件后追加字节流
 */
func FileAppendBuf(strFilePath string, bufAppend []byte) error{
	if strFilePath == ""{
		return errors.New("file path empty")
	}

	err := fileEnsureParentDir(strFilePath)
	if err != nil{
		return err
	}
	pFile,err := os.OpenFile(strFilePath,os.O_RDWR|os.O_CREATE|os.O_APPEND,FILE_PERM_DEFAULT)
	if err != nil{
		return err
	}

	_,err = pFile.Write(bufAppend)
	if err != nil{
		return err
	}

	err = pFile.Close()
	if err != nil{
		return err
	}

	return nil
}

/**
文件是否可读
 */
func FileIsAvailable(strFilePath string) bool{
	_,err := os.Lstat(strFilePath)
	return err == nil
}

/**
把文件按行读取到字符串数组
 */
func FileReadLinesToArr(strFilePath string) []string{
	if strFilePath == ""{
		return []string{}
	}

	arrLines := make([]string,0)

	fpFrom,err := os.Open(strFilePath)
	if err != nil{
		return []string{}
	}
	pScanner := bufio.NewScanner(fpFrom)
	for pScanner.Scan() {
		strLine := pScanner.Text()
		arrLines = append(arrLines,strLine)
	}
	fpFrom.Close()

	return arrLines
}

/**
把文件内容读取为字节流
 */
func FileReadBuf(strFilePath string) []byte{
	if strFilePath == ""{
		return nil
	}

	bufResult,err := ioutil.ReadFile(strFilePath)
	if err != nil{
		return nil
	}
	return bufResult
}


/**
把目录下的文件或目录读到一个数组，如果是子目录，则会在后面加上/
will add a slash to every dir
 */
func FileBrowseDirToArrName(strDir string) []string{
	arrName := make([]string,0)

	arrFileInfo,err := ioutil.ReadDir(strDir)
	if err != nil{
		return arrName
	}

	for i := 0; i < len(arrFileInfo); i++{
		tempFileInfo := arrFileInfo[i]
		if tempFileInfo.IsDir(){
			arrName = append(arrName,tempFileInfo.Name() + "/")
		}else{
			arrName = append(arrName,tempFileInfo.Name())
		}
	}

	return arrName
}


/**
把目录下的文件或目录读到一个数组，此数组内的元素都是绝对路径，如果是子目录，则会在后面加上/
will add a slash to every dir
*/
func FileBrowseDirToArrAbs(strDir string) []string{
	arrName := make([]string,0)

	FileBrowseDirToArrAbsRecur(strDir,arrName)

	return arrName
}

/**
把目录下的文件或目录读到一个数组，遍历所有子目录，如果是子目录，则会在后面加上/
will add a slash to every dir
*/
func FileBrowseDirToArrAbsRecur(strDir string, arrAbs []string){

	arrFileInfo,err := ioutil.ReadDir(strDir)
	if err != nil{
		return
	}

	strFixedDir := strDir
	if !strings.HasSuffix(strDir,"/"){
		strFixedDir += "/"
	}

	arrAbs = append(arrAbs,strFixedDir)

	for i := 0; i < len(arrFileInfo); i++{
		tempFileInfo := arrFileInfo[i]
		if tempFileInfo.IsDir(){
			FileBrowseDirToArrAbsRecur(strFixedDir + tempFileInfo.Name(),arrAbs)
		}else{
			arrAbs = append(arrAbs,strFixedDir + tempFileInfo.Name())
		}
	}

}

/**
删除一个路径
 */
func FilePathDelete(strPath string) error{
	err := os.RemoveAll(strPath)
	return err
}

/**
创建目录，会创建多级目录
 */
func FileMakeDirAll(strDir string) error{
	err := os.MkdirAll(strDir, FILE_PERM_DEFAULT)
	return err
}

/**
拷贝整个目录
 */
func FileCopyDir(strDirSrc string, strDirDest string, strCopyAction string) error{
	if strDirSrc == strDirDest {
		return errors.New("FileCopyDirOneLevel:src dir equals dest dir")
	}

	strFixedSrcDir := strDirSrc
	if !strings.HasSuffix(strFixedSrcDir,"/"){
		strFixedSrcDir += "/"
	}

	strFixedDestDir := strDirDest
	if !strings.HasSuffix(strFixedDestDir,"/"){
		strFixedDestDir += "/"
	}

	if strings.HasPrefix(strFixedDestDir,strFixedSrcDir) {
		return errors.New("FileCopyDirOneLevel:src dir covers dest dir")
	}

	return FileCopyDirRecur(strDirSrc,strDirDest,strCopyAction)
}

/**
仅仅拷贝单层目录
 */
func FileCopyDirOneLevel(strDirSrc string, strDirDest string, strCopyAction string) error{

	if strDirSrc == strDirDest {
		return errors.New("FileCopyDirOneLevel:src dir equals dest dir")
	}

	strFixedSrcDir := strDirSrc
	if !strings.HasSuffix(strFixedSrcDir,"/"){
		strFixedSrcDir += "/"
	}

	strFixedDestDir := strDirDest
	if !strings.HasSuffix(strFixedDestDir,"/"){
		strFixedDestDir += "/"
	}

	if strings.HasPrefix(strFixedDestDir,strFixedSrcDir) {
		return errors.New("FileCopyDirOneLevel:src dir covers dest dir")
	}

	arrFileInfo,err := ioutil.ReadDir(strDirSrc)
	if err != nil{
		return err
	}

	if !FileIsAvailable(strDirDest){
		err := os.MkdirAll(strDirDest,FILE_PERM_DEFAULT)
		if err != nil{
			return err
		}
	}

	for i := 0; i < len(arrFileInfo); i++{
		tempFileInfo := arrFileInfo[i]
		if tempFileInfo.IsDir(){
			err := os.MkdirAll(strFixedDestDir + tempFileInfo.Name(),FILE_PERM_DEFAULT)
			if err != nil{
				return err
			}
		}else{
			if strCopyAction == FILE_COPY_ACTION_SKIP_EXISTED{
				if !FileIsAvailable(strFixedDestDir + tempFileInfo.Name()){
					err := FileCopyBytes(strFixedSrcDir + tempFileInfo.Name(),strFixedDestDir + tempFileInfo.Name())
					if err != nil{
						return err
					}
				}
			}else if strCopyAction == FILE_COPY_ACTION_OVERWRITE_EXISTED{
				err := FileCopyBytes(strFixedSrcDir + tempFileInfo.Name(),strFixedDestDir + tempFileInfo.Name())
				if err != nil{
					return err
				}
			}
		}
	}

	return nil
}


func FileCopyDirRecur(strDirSrc string, strDirDest string, strCopyAction string) error{
	arrFileInfo,err := ioutil.ReadDir(strDirSrc)
	if err != nil{
		return err
	}

	strFixedSrcDir := strDirSrc
	if !strings.HasSuffix(strFixedSrcDir,"/"){
		strFixedSrcDir += "/"
	}

	if !FileIsAvailable(strDirDest){
		err := os.MkdirAll(strDirDest,FILE_PERM_DEFAULT)
		if err != nil{
			return err
		}
	}

	strFixedDestDir := strDirDest
	if !strings.HasSuffix(strFixedDestDir,"/"){
		strFixedDestDir += "/"
	}

	for i := 0; i < len(arrFileInfo); i++{
		tempFileInfo := arrFileInfo[i]
		if tempFileInfo.IsDir(){
			err := FileCopyDirRecur(strFixedSrcDir + tempFileInfo.Name(),strFixedDestDir + tempFileInfo.Name(),strCopyAction )
			if err != nil{
				return err
			}
		}else{
			if strCopyAction == FILE_COPY_ACTION_SKIP_EXISTED{
				if !FileIsAvailable(strFixedDestDir + tempFileInfo.Name()){
					err := FileCopyBytes(strFixedSrcDir + tempFileInfo.Name(),strFixedDestDir + tempFileInfo.Name())
					if err != nil{
						return err
					}
				}
			}else if strCopyAction == FILE_COPY_ACTION_OVERWRITE_EXISTED{
				err := FileCopyBytes(strFixedSrcDir + tempFileInfo.Name(),strFixedDestDir + tempFileInfo.Name())
				if err != nil{
					return err
				}
			}
		}
	}

	return nil
}

/**
拷贝文件
 */
func FileCopyBytes(strSrcPath string, strDestPath string) error{
	if strSrcPath == "" || strDestPath == ""{
		return errors.New("FileCopyBytes:file path empty")
	}

	iSrcSize := FileSize(strSrcPath)
	var iReadTotal int64 = 0
	bufTemp := make([]byte,FILE_TEMP_BUF_DEFAULT)
	handleSrc,err := os.Open(strSrcPath)
	if err != nil{
		return nil
	}
	defer func() {
		if handleSrc != nil{
			handleSrc.Close()
		}
	}()

	handleDest,err := os.OpenFile(strDestPath,os.O_RDWR|os.O_CREATE,FILE_PERM_DEFAULT)
	if err != nil{
		return nil
	}
	err = handleDest.Truncate(0)
	if err != nil{
		return err
	}
	defer func() {
		if handleDest != nil{
			handleDest.Close()
		}
	}()

	if iSrcSize == 0{
		return nil
	}
	for ;; {

		iReadTemp,err := handleSrc.Read(bufTemp)
		if err == nil || err == io.EOF{
			if iReadTemp > 0{
				if iReadTemp == len(bufTemp){
					_,err = handleDest.Write(bufTemp)
					if err != nil{
						return err
					}
				}else{
					_,err = handleDest.Write(bufTemp[0:iReadTemp])
					if err != nil{
						return err
					}
				}
				iReadTotal += int64(iReadTemp)
			}
		}

		if err != nil && err != io.EOF{
			return err
		}

		if iReadTotal >= iSrcSize{
			break
		}
	}

	return nil
}

/**
分段读取文件
 */
func FileReadPart(strPath string, iFileOffset int64, iLen int, bufDest []byte, iBufOffset int) error{
	pFile,err := os.Open(strPath)
	if err != nil{
		return err
	}

	if bufDest == nil || len(bufDest) < int(iLen){
		return errors.New("FileReadPart:buf size is less then required")
	}

	if iFileOffset > 0{
		iRealOffset,err := pFile.Seek(iFileOffset,0)
		if err != nil{
			pFile.Close()
			return err
		}

		if iRealOffset != iFileOffset{
			return errors.New("FileReadPart:offset error")
		}
	}

	bufTemp := make([]byte,FILE_TEMP_BUF_DEFAULT)
	iReadTotal := 0

	//pReader := bufio.NewReader(pFile)
	for ; ; {
		iReadTemp,errRead := pFile.Read(bufTemp)
		if errRead != nil && errRead != io.EOF{
			return errRead
		}

		iWriteTemp := iReadTemp
		if iReadTemp + iReadTotal > iLen{
			iWriteTemp = iLen - iReadTotal
		}

		ByteCopy(bufTemp,0, bufDest,iBufOffset + iReadTotal,iWriteTemp)
		iReadTotal += iWriteTemp

		if errRead != nil || iReadTotal >= iLen{
			break
		}
	}

	if bufDest != nil{
		iWritten,err := pFile.Write(bufDest)
		if err != nil{
			pFile.Close()
			return err
		}

		if iWritten != len(bufDest){
			return errors.New("bytes written error")
		}
	}

	pFile.Close()
	return nil
}

/**
分段写入文件
 */
func FileWritePart(strPath string, iFileOffset int64, bufToWrite []byte) error{
	pFile,err := os.OpenFile(strPath,os.O_RDWR|os.O_CREATE,FILE_PERM_DEFAULT)
	if err != nil{
		return err
	}

	if iFileOffset > 0{
		iRealOffset,err := pFile.Seek(iFileOffset,0)
		if err != nil{
			pFile.Close()
			return err
		}

		if iRealOffset != iFileOffset {
			return errors.New("offset error")
		}
	}
	if bufToWrite != nil{
		iWritten,err := pFile.Write(bufToWrite)
		if err != nil{
			pFile.Close()
			return err
		}

		if iWritten != len(bufToWrite){
			return errors.New("bytes written error")
		}
	}

	pFile.Close()
	return nil
}

/**
设置文件大小
 */
func FileSetSize(strPath string, iOffset int64) error{
	pFile,err := os.OpenFile(strPath,os.O_RDWR|os.O_CREATE,FILE_PERM_DEFAULT)
	if err != nil{
		return err
	}
	if iOffset > 0{
		err = pFile.Truncate(iOffset)
		if err != nil{
			pFile.Close()
			return err
		}
	}

	pFile.Close()
	return nil
}

/**
清空一个目录
 */
func FileClearSubDir(strDir string) error{
	arrFileInfo,err := ioutil.ReadDir(strDir)
	if err != nil{
		return err
	}

	strFixedSrcDir := strDir
	if !strings.HasSuffix(strFixedSrcDir,"/"){
		strFixedSrcDir += "/"
	}

	for i := 0; i < len(arrFileInfo); i++{
		tempFileInfo := arrFileInfo[i]
		if tempFileInfo.IsDir(){
			err := os.RemoveAll(strFixedSrcDir + tempFileInfo.Name())
			if err != nil{
				return err
			}
		}else{
			err := os.Remove(strFixedSrcDir + tempFileInfo.Name())
			if err != nil{
				return err
			}
		}
	}

	return nil
}

/**
移除一个目录
 */
func FileRemoveDir(strDir string) error{
	err := os.RemoveAll(strDir)
	return err
}

/**
移除一个文件
 */
func FileRemoveFile(strDir string) error{
	err := os.Remove(strDir)
	return err
}

func FileGetUtfBom() []byte{
	return FILE_BOM_UTF8
}

func FileGetUnicodeBom() []byte{
	return FILE_BOM_UNICODE
}


