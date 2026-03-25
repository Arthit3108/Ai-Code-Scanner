pipeline {
    agent any

    environment {
        // Create a 'Secret text' credential in Jenkins with ID 'gemini-api-key'
        GEMINI_API_KEY = credentials('gemini-api-key')
        SCAN_TARGET = '.'
        SCAN_SEVERITY = 'HIGH,CRITICAL'
        OUTPUT_TYPE = 'json' // or 'sarif'
    }

    stages {
        stage('Install Prerequisites') {
            steps {
                sh '''
                    # Install Trivy if not present
                    if ! command -v trivy &> /dev/null; then
                        echo "Installing Trivy..."
                        curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin
                    fi
                    
                    # Ensure jq is installed for threshold checking
                    if ! command -v jq &> /dev/null; then
                        echo "Installing jq..."
                        if [ -f /etc/debian_version ]; then
                            sudo apt-get update && sudo apt-get install -y jq
                        elif [ -f /etc/redhat-release ]; then
                            sudo yum install -y jq
                        fi
                    fi
                '''
            }
        }

        stage('Download Scanner') {
            steps {
                sh '''
                    VERSION="v1" # Specify version or use 'latest'
                    echo "Downloading AI Code Scanner ${VERSION}..."
                    curl -sfL "https://github.com/Arthit3108/ai-code-scanner/releases/download/${VERSION}/ai-scanner" -o ai-scanner
                    chmod +x ai-scanner
                '''
            }
        }

        stage('AI Security Scan') {
            steps {
                sh '''
                    ./ai-scanner scan \
                        --target "${SCAN_TARGET}" \
                        --severity "${SCAN_SEVERITY}" \
                        --output-type "${OUTPUT_TYPE}" \
                        --output "scanResult.${OUTPUT_TYPE}"
                '''
            }
        }

        stage('Check Severity Thresholds') {
            when {
                expression { env.OUTPUT_TYPE == 'json' }
            }
            steps {
                sh '''
                    CRITICAL_COUNT=$(jq '[.vulnerabilities[]? | select(.severity == "CRITICAL")] | length' scanResult.json)
                    HIGH_COUNT=$(jq '[.vulnerabilities[]? | select(.severity == "HIGH")] | length' scanResult.json)
                    
                    echo "---------------------------------------"
                    echo "Scan Results Summary:"
                    echo "Critical Vulnerabilities: $CRITICAL_COUNT"
                    echo "High Vulnerabilities: $HIGH_COUNT"
                    echo "---------------------------------------"

                    if [ "$CRITICAL_COUNT" -gt 0 ] || [ "$HIGH_COUNT" -gt 0 ]; then
                        echo "❌ Security check failed: High/Critical vulnerabilities found."
                        exit 1
                    fi
                '''
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: "scanResult.*", fingerprint: true
        }
    }
}
