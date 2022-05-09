const core = require('@actions/core');
const { promises: fs } = require('fs')
const funcs = require('./bump')

const actions = {
  major: funcs.incMajor,
  minor: funcs.incMinor,
  patch: funcs.incPatch,
}

async function run() {
  try {
    const preReleaseLabel = core.getInput('pre-release');
    const releaseType = core.getInput('release-type');

    if (!(releaseType in actions)) {
      throw new Error("unexpected value for release types.")
    }

    const file = core.getInput('path');
    let content = await fs.readFile(file, 'utf8')
    let currentVersion = content.trim()
    if (! await funcs.isValidSemver(currentVersion)) {
      throw new Error('File did not contain valid semver')
    }

    let oldVersion = currentVersion
    let newVersion = await actions[releaseType](currentVersion, preReleaseLabel)
    if (newVersion === null) {
      throw new Error("unable to update version")
    }

    core.setOutput('previous-version', oldVersion);
    core.setOutput('version', newVersion);
  } catch (error) {
    core.setFailed(error.message);
  }
}

run();
