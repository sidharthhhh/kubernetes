# Setting up the AWS CLI

To interact with AWS services directly from your terminal, you need the AWS CLI (Command Line Interface). This guide will help you install it and link it securely to your AWS account.

## Step 1: Install the AWS CLI

### Windows
1. Download the official installer: [AWS CLI MSI Installer](https://awscli.amazonaws.com/AWSCLIV2.msi)
2. Run the downloaded file and follow the on-screen prompts.
3. Open a new PowerShell terminal and verify:
   ```powershell
   aws --version
   ```

### macOS
1. Open your terminal and run:
   ```bash
   curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "AWSCLIV2.pkg"
   sudo installer -pkg AWSCLIV2.pkg -target /
   ```
2. Verify:
   ```bash
   aws --version
   ```

### Linux
1. Open your terminal and run:
   ```bash
   curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
   unzip awscliv2.zip
   sudo ./aws/install
   ```
2. Verify:
   ```bash
   aws --version
   ```

---

## Step 2: Create an IAM User
You should *never* use your root AWS account for daily tasks. Instead, we'll create an IAM User with Programmatic Access.

1. Go to the AWS Management Console and open **IAM**.
2. Click **Users** on the left menu, then **Create user**.
3. Enter a username (e.g., `kubernetes-admin`).
4. We only need CLI access, so you do *not* need to check "Provide user access to the AWS Management Console". Click **Next**.
5. Select **Attach policies directly**.
6. Search for and select **AdministratorAccess**. *(Note: This gives full permissions. In a strict production environment, you would lock this down further, but it's required for learning EKS).*
7. Click **Next**, then **Create user**.

---

## Step 3: Generate Access Keys
1. Click on your newly created user (`kubernetes-admin`).
2. Go to the **Security credentials** tab.
3. Scroll down to **Access keys** and click **Create access key**.
4. Select **Command Line Interface (CLI)**.
5. Click **Next** and **Create access key**.
6. **IMPORTANT:** Do not close this screen yet! You need the **Access key ID** and **Secret access key** for the next step.

---

## Step 4: Configure the CLI
Go back to your local terminal and run:

```bash
aws configure
```

It will prompt you for four things:
- **AWS Access Key ID**: Paste the Access Key ID from Step 3.
- **AWS Secret Access Key**: Paste the Secret Access Key from Step 3.
- **Default region name**: Type your preferred region (e.g., `us-east-1`, `eu-west-1`, or `ap-south-1`).
- **Default output format**: Type `json`.

**That's it!** Your terminal is now authenticated and ready to interact with AWS. You can test it by running:
```bash
aws sts get-caller-identity
```
