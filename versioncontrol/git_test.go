const GIT_NAME string = "Git"
	result := VCMap[GIT_NAME].NameOfDir()
	result := VCMap[GIT_NAME].NameOfVC()
	VCMap[GIT_NAME].moveToDir(t)
	diff, err := VCMap[GIT_NAME].GetDiff()
	VCMap[GIT_NAME].moveToDir(t)
	err := VCMap[GIT_NAME].SetHooks(HomeDir)
		filePath := filepath.Join(VCMap[GIT_NAME].NameOfDir(), "hooks", fileName)
