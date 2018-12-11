const fs = require("fs");
const { promisify } = require("util");
const snapshotDiff = require("snapshot-diff");

const readFile = promisify(fs.readFile);

async function getEnvYaml(env) {
  return [
    await readFile(`../examples/${env}/deployment.yml`, "utf8"),
    env !== "feature"
      ? await readFile(`../examples/${env}/ingress.yml`, "utf8")
      : "",
    await readFile(`../examples/${env}/service.yml`, "utf8")
  ].join("\n---\n");
}

test("feature <> staging", async () => {
  const feature = await getEnvYaml("feature");
  const staging = await getEnvYaml("staging");

  expect(
    snapshotDiff(feature, staging, {
      contextLines: 0,
      aAnnotation: "feature",
      bAnnotation: "staging"
    })
  ).toMatchSnapshot();
});

test("staging <> production", async () => {
  const staging = await getEnvYaml("staging");
  const production = await getEnvYaml("production");

  expect(
    snapshotDiff(staging, production, {
      contextLines: 0,
      aAnnotation: "staging",
      bAnnotation: "production"
    })
  ).toMatchSnapshot();
});
