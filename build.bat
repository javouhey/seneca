@ECHO OFF

SET _VER=
SET _OLDGOBIN=
SET _OLDGOPATH=
SET _GOPATH=0 
set _GOBIN=0
set _=
set _LDFLAGS=

SET _B=src\github.com\javouhey
SET _A=src\github.com\javouhey\seneca

IF exist %_A%\ ( 
  echo.src directory found. Good 
  goto next

) ELSE ( 
  echo."WARNING: Creating symbolic link on windows requires administrator privilege"
  mkdir %_B% && echo.src created
  cd %_B%
  mklink /D seneca ..\..\.. 
  if errorlevel 1 (
      cd ..\..\..
      goto die
  )
  echo.do
  cd ..\..\..
  echo.done creating src
  goto next
)

:next

REM echo."%GOBIN%"
REM echo."%GOPATH%"

FOR /F "usebackq" %%i IN (`type version`) DO SET _VER=%%i
SET _=%CD%
FOR /F "usebackq" %%i IN (`git rev-parse HEAD`) DO SET _GITSHA=%%i

IF /I "%GOBIN%"=="" ( ECHO.ignore GOBIN) ELSE ( 
  ECHO.GOBIN env found 
  set _GOBIN=1
  SET _OLDGOBIN=%GOBIN%
)
IF /I "%GOPATH%"=="" ( ECHO.ignore GOPATH) ELSE ( 
  ECHO.GOPATH env found
  set _GOPATH=1
  SET _OLDGOPATH=%GOPATH%
)

REM change gopath to local dir
set GOPATH=%_%
REM  -- set GOBIN=%_%\bin

set _LDFLAGS="-X main.GitSHA '%_GITSHA%' -X main.Version '%_VER%' -w"

echo.Installing seneca to %_%\bin
go install -ldflags %_LDFLAGS% github.com/javouhey/seneca


REM echo.%_OLDGOPATH%
REM echo.%_OLDGOBIN%

IF %_GOPATH%==1 (
  echo.RESTORE GOPATH
  set GOPATH=%_OLDGOPATH%
)
IF %_GOBIN%==1 (
  echo.RESTORE GOBIN
  set GOBIN=%_OLDGOBIN%
)
goto end

:die
echo.Aborting

:end
echo.Bye
