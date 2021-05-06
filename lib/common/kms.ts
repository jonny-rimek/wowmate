import cdk = require('@aws-cdk/core');
import * as kms from "@aws-cdk/aws-kms";

interface Props extends cdk.StackProps {
}

export class Kms extends cdk.Construct {
    public readonly key: kms.Key

    constructor(scope: cdk.Construct, id: string, props?: Props) {
        super(scope, id)

        this.key = new kms.Key(this, 'WowmateKey', {
            enableKeyRotation: true,
        })
    }
}
