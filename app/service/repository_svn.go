package service

import (
	"errors"
	"github.com/zhenzhongfu/gopub/app/libs"
	"path/filepath"
	_ "github.com/astaxie/beego"
	_ "strconv"
	_ "strings"
)

type repositoryServiceSvn struct{}

// 返回一个仓库对象
func (this *repositoryServiceSvn) GetRepoByProjectId(projectId int) (*RepositorySvn, error) {
	project, err := ProjectService.GetProject(projectId)
	if err != nil {
		return nil, err
	}
	return OpenRepositorySvn(project.Domain)
}

// 克隆svn仓库
func (this *repositoryServiceSvn) CloneRepo(url string, dst string) error {
	out, stderr, err := libs.ExecCmd("svn", "co", url, dst)
	debug("out", out)
	debug("stderr", stderr)
	debug("err", err)
	if err != nil {
		return concatenateError(err, stderr)
	}
	return nil
}

type RepositorySvn struct {
	Path string
}

func OpenRepositorySvn(repoPath string) (*RepositorySvn, error) {
	repoPath = GetProjectPath(repoPath)
	repoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	} else if !libs.IsDir(repoPath) {
		return nil, errors.New("no such file or directory")
	}

	return &RepositorySvn{Path: repoPath}, nil
}

// 导出版本到tar包
func (repo *RepositorySvn) Export(startVer, endVer string, filename string) error {
	//svn export -r 15671 repo.Path tmp|awk '{print $2}'|head -n 1|xargs tar zcvf filename
	cmd := "svn export -r " + endVer + " " + repo.Path + " --force |awk {'print $2'} > tmp; sleep 60; head -n 1 tmp|xargs tar zcvf " + filename
	_, stderr, err := libs.ExecCmdDir(repo.Path, "/bin/bash", "-c", cmd)

	if err != nil {
		return concatenateError(err, stderr)
	}
	return nil
}
