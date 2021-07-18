{-# LANGUAGE CPP #-}
{-# LANGUAGE NoRebindableSyntax #-}
{-# OPTIONS_GHC -fno-warn-missing-import-lists #-}
module Paths_hs_git_clone (
    version,
    getBinDir, getLibDir, getDynLibDir, getDataDir, getLibexecDir,
    getDataFileName, getSysconfDir
  ) where

import qualified Control.Exception as Exception
import Data.Version (Version(..))
import System.Environment (getEnv)
import Prelude

#if defined(VERSION_base)

#if MIN_VERSION_base(4,0,0)
catchIO :: IO a -> (Exception.IOException -> IO a) -> IO a
#else
catchIO :: IO a -> (Exception.Exception -> IO a) -> IO a
#endif

#else
catchIO :: IO a -> (Exception.IOException -> IO a) -> IO a
#endif
catchIO = Exception.catch

version :: Version
version = Version [0,1,0,0] []
bindir, libdir, dynlibdir, datadir, libexecdir, sysconfdir :: FilePath

bindir     = "/Users/rohitpaulk/experiments/codecrafters/testers/tester-utils/test_helpers/executable_test/haskell/.stack-work/install/x86_64-osx/00a12b6dee8c353e2c7135411b34d6ff22efe969da779bbc615e68ee22b36a6c/8.8.3/bin"
libdir     = "/Users/rohitpaulk/experiments/codecrafters/testers/tester-utils/test_helpers/executable_test/haskell/.stack-work/install/x86_64-osx/00a12b6dee8c353e2c7135411b34d6ff22efe969da779bbc615e68ee22b36a6c/8.8.3/lib/x86_64-osx-ghc-8.8.3/hs-git-clone-0.1.0.0-CmkJQUbgsbX5siKfNWiu56-hs-git-clone-exe"
dynlibdir  = "/Users/rohitpaulk/experiments/codecrafters/testers/tester-utils/test_helpers/executable_test/haskell/.stack-work/install/x86_64-osx/00a12b6dee8c353e2c7135411b34d6ff22efe969da779bbc615e68ee22b36a6c/8.8.3/lib/x86_64-osx-ghc-8.8.3"
datadir    = "/Users/rohitpaulk/experiments/codecrafters/testers/tester-utils/test_helpers/executable_test/haskell/.stack-work/install/x86_64-osx/00a12b6dee8c353e2c7135411b34d6ff22efe969da779bbc615e68ee22b36a6c/8.8.3/share/x86_64-osx-ghc-8.8.3/hs-git-clone-0.1.0.0"
libexecdir = "/Users/rohitpaulk/experiments/codecrafters/testers/tester-utils/test_helpers/executable_test/haskell/.stack-work/install/x86_64-osx/00a12b6dee8c353e2c7135411b34d6ff22efe969da779bbc615e68ee22b36a6c/8.8.3/libexec/x86_64-osx-ghc-8.8.3/hs-git-clone-0.1.0.0"
sysconfdir = "/Users/rohitpaulk/experiments/codecrafters/testers/tester-utils/test_helpers/executable_test/haskell/.stack-work/install/x86_64-osx/00a12b6dee8c353e2c7135411b34d6ff22efe969da779bbc615e68ee22b36a6c/8.8.3/etc"

getBinDir, getLibDir, getDynLibDir, getDataDir, getLibexecDir, getSysconfDir :: IO FilePath
getBinDir = catchIO (getEnv "hs_git_clone_bindir") (\_ -> return bindir)
getLibDir = catchIO (getEnv "hs_git_clone_libdir") (\_ -> return libdir)
getDynLibDir = catchIO (getEnv "hs_git_clone_dynlibdir") (\_ -> return dynlibdir)
getDataDir = catchIO (getEnv "hs_git_clone_datadir") (\_ -> return datadir)
getLibexecDir = catchIO (getEnv "hs_git_clone_libexecdir") (\_ -> return libexecdir)
getSysconfDir = catchIO (getEnv "hs_git_clone_sysconfdir") (\_ -> return sysconfdir)

getDataFileName :: FilePath -> IO FilePath
getDataFileName name = do
  dir <- getDataDir
  return (dir ++ "/" ++ name)
