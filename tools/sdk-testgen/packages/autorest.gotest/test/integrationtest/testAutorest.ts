import * as assert from 'assert';
import * as fs from 'fs';
import * as path from 'path';
import { Helper } from '@autorest/testmodeler/dist/src/util/helper';
import { exec } from 'child_process';

async function compare(dir1: string, dir2: string) {
    const cmd = 'diff -r --exclude=gen.zip --strip-trailing-cr -I _filePath -I x-ms-original-file -I file:/// ' + dir1 + ' ' + dir2;
    console.log(cmd);
    return await new Promise<boolean>((resolve, reject) => {
        exec(cmd, (error, stdout) => {
            if (error) {
                console.log('exec error: ' + error + ', ' + stdout);
                // Reject if there is an error:
                return reject(false);
            }
            // Otherwise resolve the promise:
            return resolve(true);
        });
    });
}

async function runAutorest(readmePath: string, extraOption: string[]) {
    const cmd =
        path.join(`${__dirname}`, '..', '..' + '/node_modules', '.bin', 'autorest') +
        ' --version=3.7.3 --testmodeler.generate-mock-test --testmodeler.generate-sdk-example --testmodeler.generate-scenario-test --testmodeler.generate-sdk-sample --use=' +
        path.join(`${__dirname}`, '..', '..', '..', 'autorest.testmodeler') +
        ' --use=' +
        path.join(`${__dirname}`, '..', '..') +
        ' ' +
        readmePath +
        ' ' +
        extraOption.join(' ');
    console.log(cmd);
    return await new Promise<boolean>((resolve, reject) => {
        exec(cmd, (error, stdout, stderr) => {
            if (error) {
                console.error('exec error: ' + stderr);
                // Reject if there is an error:
                return reject(false);
            }
            // Otherwise resolve the promise:
            console.log(stdout);
            return resolve(true);
        });
    });
}
function getExtraOption(outputFolder: string) {
    return [`--output-folder=${outputFolder}`, '--use=@autorest/go@latest', '--file-prefix="zz_generated_"', '--track2', '--go', '--debug', '--module-version=0.1.0'];
}

async function runSingleTest(swaggerDir: string, rp: string, extraOption: string[], outputFolder: string, tempOutputFolder: string): Promise<boolean> {
    let result = true;
    let msg = '';
    const readmePath = path.join(swaggerDir, 'specification', rp, 'resource-manager/readme.md');
    await runAutorest(readmePath, extraOption)
        .then((res) => {
            if (!res) {
                msg = 'Run autorest not successfully!';
            }
            result = res;
        })
        .catch((e) => {
            msg = `Run autorest failed! ${e}`;
            result = false;
        });
    if (result) {
        await compare(outputFolder, tempOutputFolder)
            .then((res1) => {
                if (res1 === false) {
                    msg = 'The generated files have changed!';
                }
                result = res1;
            })
            .catch((e) => {
                msg = `The diff has some error ${e}`;
                result = false;
            });
    } else {
        console.error(msg);
    }
    return result;
}

describe('Run autorest and compare the output', () => {
    beforeAll(async () => {
        //
    });

    afterAll(async () => {
        //
    });

    it('', async () => {
        jest.setTimeout(8 * 60 * 60 * 1000);
        const swaggerDir = path.join(`${__dirname}`, '../../../../swagger');
        const outputDir = path.join(`${__dirname}`, 'output');
        const tempoutputDir = path.join(`${__dirname}`, 'tempoutput');

        let finalResult = true;
        const allTests: Array<Promise<boolean>> = [];
        for (const rp of ['appplatform', 'compute', 'signalr']) {
            console.log('Start Processing: ' + rp);

            // Remove tmpoutput
            const outputFolder = path.join(outputDir, rp, 'test');
            const tempOutputFolder = path.join(tempoutputDir, rp, 'test');
            Helper.deleteFolderRecursive(tempOutputFolder);
            fs.mkdirSync(tempOutputFolder, { recursive: true });

            const test = runSingleTest(swaggerDir, rp, getExtraOption(tempOutputFolder), outputFolder, tempOutputFolder);
            allTests.push(test);
        }
        if ((process.env['PARALELL_TEST'] || 'false').toLowerCase() === 'true') {
            finalResult = (await Promise.all(allTests)).every((x) => x);
        } else {
            for (const test of allTests) {
                if (!(await test)) {
                    assert.fail('Test failed!');
                }
            }
        }

        assert.strictEqual(finalResult, true);
    });
});
