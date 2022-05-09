const semver = require("semver");

let isValidSemver = function (version) {
  return new Promise((resolve)=>{
    resolve(semver.valid(version) !== null)
  })
}

let incMajor = function (version, preReleaseLabel) {
  return new Promise((resolve) => {
    let type = 'major'
    if (preReleaseLabel) {
      type = 'premajor'
    }
    resolve(semver.inc(version, type, preReleaseLabel))
  })
}

let incMinor = function (version, preReleaseLabel) {
  return new Promise((resolve) => {
    let type = 'minor'
    if (preReleaseLabel) {
      type = 'preminor'
    }
    resolve(semver.inc(version, type, preReleaseLabel))
  })
}

let incPatch = function (version, preReleaseLabel) {
  return new Promise((resolve) => {
    let type = 'patch'
    if (preReleaseLabel) {
      type = 'prepatch'
    }
    resolve(semver.inc(version, type, preReleaseLabel))
  })
}


module.exports = {incMajor, incMinor, incPatch, isValidSemver};
