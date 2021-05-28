package login_sdk_go

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestValidateHmacToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	loginSgk := New(Options{
		ShaSecretKey: "your-256-bit-secret",
	})
	_, err := loginSgk.Validate(token)

	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateRsaToken(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6InNnRnk0NjRrVk5YVFo2YmVYM0tFT2kyam1yWnA4bUQiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOltdLCJlbWFpbCI6ImxvZ2luLXNka0B0ZXN0LmNvbSIsImV4cCI6MTYyMjEyMjkyMiwiZ3JvdXBzIjpbXSwiaWF0IjoxNjIyMTE5MzIyLCJpc19tYXN0ZXIiOnRydWUsImlzcyI6Imh0dHBzOi8vbG9naW4ueHNvbGxhLmNvbSIsImp0aSI6Ijc2YmNjYWRkLWMwNDEtNGFmOS1hN2QzLWFhNDEyNjc3ZjU5OCIsInByb21vX2VtYWlsX2FncmVlbWVudCI6dHJ1ZSwic2NwIjpbIm9mZmxpbmUiXSwic3ViIjoiOTNlMWE5YTMtOGRkZC00OTMxLWE4ZmQtMjA5MGY2M2VkZmI4IiwidHlwZSI6Inhzb2xsYV9sb2dpbiIsInVzZXJuYW1lIjoibG9naW4tc2RrLXRlc3QiLCJ4c29sbGFfbG9naW5fYWNjZXNzX2tleSI6ImZPOGNPc2o4TjdHSWw3MldscVNsYndnSmk0ZGRBOHVQbFd2ZDYxTk1ydnMiLCJ4c29sbGFfbG9naW5fcHJvamVjdF9pZCI6IjQwZGIyZWE0LTVkNDItMTFlNi1hM2ZmLTAwNTA1NmEwZTA0YSJ9.Toe6NdnnXdCtrVOodRotmpUoU79c4FDlRVi01vIoQNaCICqDeTr54SVYpceAm0xTc_f7MlyJVw9pWodtBcN5Ehq9cTARwIr_7iTx_QN_EdjA6twUneCFq3FJBgVvXmCp2-foG1r-rUl0GLrv_C2DOf8e24LfQ7gtBKFhrg-wwuPydO9zUSmhs7qgM-vMRjiXIM8fx-YYVZB1jB7Ik8hU89dWpsbrY4C4MR8kVen32V-uOVDUCJ1Ao6pG8U7RyWCrX3DZiiQDmg1_vCeAseY-VIyI5-Ta30FsC42r5jPbQvUeXfLQKOfRwiJD-5RpEau-Dz7C2BVbtIImRIXHiHodpA"
	loginSgk := New(Options{
		LoginProjectId: "40db2ea4-5d42-11e6-a3ff-005056a0e04a",
	})
	_, err := loginSgk.Validate(token)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRefreshToken(t *testing.T) {
	refreshToken := "Gas4JcYby1sGkd1-y6hl7rQbUZ8Cc8Z0CZcAcEIi3X8.pPpDvm0KlLK7DqSM_Zt7Ru8ODLoru62XbE4pFIIw1VM"
	clientId, _ := strconv.Atoi(os.Getenv("LOGIN_CLIENT_ID"))

	loginSgk := New(Options{
		LoginProjectId:    "40db2ea4-5d42-11e6-a3ff-005056a0e04a",
		LoginClientId:     clientId,
		LoginClientSecret: os.Getenv("LOGIN_CLIENT_SECRET"),
	})
	res, err := loginSgk.Refresh(refreshToken)

	fmt.Print(res)

	if err != nil {
		t.Fatal(err)
	}
}
