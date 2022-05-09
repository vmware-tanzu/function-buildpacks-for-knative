const bump = require('./bump');
// const process = require('process');
// const cp = require('child_process');
// const path = require('path');

test('valid semver', async () => {
  let valid = await bump.isValidSemver('0.8.9-beta.0+foo')
  expect(valid).toBe(true);
});

test('increases major', async () => {
  let newVersion = await bump.incMajor('1.2.3');
  expect(newVersion).toBe("2.0.0");
});

test('increases premajor', async () => {
  let newVersion = await bump.incMajor('1.2.3', "something");
  expect(newVersion).toBe("2.0.0-something.0");
});

test('increases minor', async () => {
  let newVersion = await bump.incMinor('1.2.3');
  expect(newVersion).toBe("1.3.0");
});

test('increases preminor', async () => {
  let newVersion = await bump.incMinor('1.2.3', "something");
  expect(newVersion).toBe("1.3.0-something.0");
});

test('increases patch', async () => {
  let newVersion = await bump.incPatch('1.2.3');
  expect(newVersion).toBe("1.2.4");
});

test('increases prepatch', async () => {
  let newVersion = await bump.incPatch('1.2.3', "something");
  expect(newVersion).toBe("1.2.4-something.0");
});

// shows how the runner will run a javascript action with env / stdout protocol
// test('test runs', () => {
//   process.env['INPUT_PRE_RELEASE'] = "minor";
//   process.env['INPUT_PATH'] = "./VERSION";
//   const ip = path.join(__dirname, 'index.js');
//   const result = cp.execSync(`node ${ip}`, {env: process.env}).toString();
//   console.log(result);
// })
