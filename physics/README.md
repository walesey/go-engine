# Installation Instructions for bullet

## Windows
    * Install Mingw to /c/mingw
    * Install Msys to /c/msys
    * Unzip msysCORE-1.0.11-bin.tar to /c/msys/1.0/
    * Unzip cmake-3.5.0-rc3-win32-x86.zip to /c/cmake
    * Add /c/cmake/bin to PATH
    * Open msys terminal
    * cd into bullet project folder.
    * cmake -G "MSYS Makefiles" -DBUILD_EXTRAS=off -DINSTALL_LIBS=on -DBUILD_SHARED_LIBS=on
    * make
    * make install
    * copy the contents of C:/Program Files (x86)/BULLET_PHYSICS/ to /c/mingw/mingw64/
    * bullet headers will now be available to mingw gcc
