const { promisify } = require('util')
const { exec } = require('child_process')
const os = require('os')
const path = require('path')

const execPromise = promisify(exec)

const isCI = !!process.env.VERCEL_ENV

if (!isCI) {
  console.log('not in a vercel CI environment, exiting...')
  return
}

async function run() {
  try {
    const nitricPath = path.join(os.homedir(), '.nitric', 'bin', 'nitric')

    // Array of commands to execute
    const commands = [
      'curl -L https://nitric.io/install?version=latest | bash',
      `mv ${nitricPath} /usr/local/bin/`,
      'nitric version',
      'nitric help > src/assets/cli-usage.txt',
    ]

    // Run each command in the array
    for (const command of commands) {
      const { stdout, stderr } = await execPromise(command)
      console.log(`Command output: ${stdout}`)
      if (stderr) {
        console.error(`Command stderr: ${stderr}`)
      }
    }
  } catch (error) {
    console.log(error.message)
  }
}

run()
