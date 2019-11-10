import { expect as expectCDK, matchTemplate, MatchStyle } from '@aws-cdk/assert';
import cdk = require('@aws-cdk/core');
import Wowmate = require('../lib/wowmate-stack');

test('Empty Stack', () => {
    const app = new cdk.App();
    // WHEN
    const stack = new Wowmate.WowmateStack(app, 'MyTestStack');
    // THEN
    expectCDK(stack).to(matchTemplate({
      "Resources": {}
    }, MatchStyle.EXACT))
});