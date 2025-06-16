# Loan Engine - Summary Requirements
## Core Business Rules
- **State Flow**: `proposed` → `approved` → `invested` → `disbursed`
- **Forward-only movement**: No backward state transitions allowed
- **Multiple investors**: Each loan can have multiple investors with individual amounts
- **Investment constraint**: Total invested amount cannot exceed loan principal
- **Auto state transition**: Loan becomes `invested` when total investment equals principal
- **Adding State Funding**: Based on my analysis, to simplify the logic before disbursement, I have introduced a new state, `Funding`, positioned between `Approved` and `Invested`.

### Workflow e2e
```mermaid
flowchart TD
sysLoan((System Loan)) -- borrower --> proposeLoan[Propose Loan]

    subgraph State: Propose
        proposeLoan -- field validator --> sb[Survey Borrower]
        sb --> inputDoc[Input Survey Document]
        inputDoc --> valDoc{Validation Document}
    end

    subgraph State: Approve
        valDoc -->|Invalid| loanReject[Loan Rejected]
        loanReject --> st(((Stop)))
        valDoc -->|Valid| approveLoan[Loan Approved]
        approveLoan --> generateTemplate[Generate Agreement Template]
    end

    subgraph State: Invest
        generateTemplate --> listingLoan[Listing Available Investment]
        listingLoan -- Investor --> makeInvestment[Make Investment]
        makeInvestment --> sendIndividualAgreement[Send Individual Agreement to Investor]
        sendIndividualAgreement --> updateTotal[Update Total Investment]
        updateTotal --> checkTotal{Total Investment = Principal?}

        checkTotal -->|No| listingLoan
        checkTotal -->|Yes| loanFullyInvested[Loan Fully Invested]
        loanFullyInvested --> sendCompletionNotification[Send Completion Notification to All Investors]
    end

    subgraph State: Disbursement
        sendCompletionNotification --> waitingDisbursement[Waiting for Disbursement]
        waitingDisbursement -- Field Officer --> collectSignedAgreement[Collect Signed Agreement from Borrower]
        collectSignedAgreement --> disbursement[Disburse to Borrower]
        disbursement --> loanActive[Loan Active]
    end

    loanActive --> e(((End)))

```
---
## Usecase Diagram
![img.png](diagram/output/usecase-diagram.png)

Based on use case diagram, below are feature we will provide.

## Feature List
### 1. Auth
**Description**: Authentication & Authorization Users

**API:**
- Login user

### 2. Loan Lifecycle Management
**Description**: Core loan workflow management from proposal to disbursement.

**API:**
- Create loan proposal
- Approve loan (PROPOSED → APPROVED)
- Disburse loan (INVESTED → DISBURSED)

---

### 3. Investment Management
**Description**: Handle multiple investor investments with aggregation.

**API:**
- List available loans for investment
- Make investment in loan
- Get investor's investment portfolio

---

### 4. Document Management
**Description**: upload documents (proofs, agreements, signed contracts).

**API:**
- Upload document files
- Download/view document

---

### 5. System Health & Monitoring
**Description**: System health checks and monitoring endpoints for operational visibility.

**API:**
- Basic health check


