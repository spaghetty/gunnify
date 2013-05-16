##
## AC macros for using go lang
##

AC_DEFUN([CHECK_GOLANG],[
        AC_ARG_VAR(GOLANG, [ Go compiler full path ])

        AC_PATH_PROG([GOLANG],[go],[no])
        if test x$GOLANG != xno; then
                found_golang="yes";
        fi
        AS_VAR_SET_IF([found_golang],
        [],
        [AC_MSG_ERROR([Cannot find golang installed. Install go and set GOROOT \
var])]
        )
])



#                                               -*- Autoconf -*-
# Process this file with autoconf to produce a configure script.

AC_PREREQ([2.69])
AC_INIT([gunnify], [0.96], [spaghetty@gmail.com])
AM_INIT_AUTOMAKE(foreign tar-ustar)
# Checks for programs.
AC_PROG_INSTALL
CHECK_GOLANG

# Checks for libraries.

# Checks for header files.

# Checks for typedefs, structures, and compiler characteristics.

# Checks for library functions.

AC_CONFIG_FILES([Makefile])
AC_OUTPUT
