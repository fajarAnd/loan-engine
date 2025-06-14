# Loan Engine - 📋 Summary Requirements
## Core Business Rules
- **State Flow**: `proposed` → `approved` → `invested` → `disbursed`
- **Forward-only movement**: No backward state transitions allowed
- **Multiple investors**: Each loan can have multiple investors with individual amounts
- **Investment constraint**: Total invested amount cannot exceed loan principal
- **Auto state transition**: Loan becomes `invested` when total investment equals principal

### Workflow e2e
```mermaid
flowchart TD
sysLoan((System Loan)) -- borrower --> proposeLoan[Propose Loan]

    subgraph State: Propose
        proposeLoan -- field officer --> sb[Survey Borrower]
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

### 1. Loan Lifecycle Management
**Description**: Core loan workflow management from proposal to disbursement.

**API:**
- `POST /api/v1/loans` - Create loan proposal
- `GET /api/v1/loans/{id}` - Get loan details
- `PUT /api/v1/loans/{id}/approve` - Approve loan (PROPOSED → APPROVED)
- `PUT /api/v1/loans/{id}/disburse` - Disburse loan (INVESTED → DISBURSED)

---

### 2. Investment Management
**Description**: Handle multiple investor investments with real-time aggregation and automatic state transitions.

**API:**
- `GET /api/v1/loans/available` - List available loans for investment
- `POST /api/v1/loans/{id}/investments` - Make investment in loan
- `GET /api/v1/investors/{investor_id}/portfolio` - Get investor's investment portfolio

---

### 3. Document Management
**Description**: Secure file upload and management for loan-related documents (proofs, agreements, signed contracts).

**API:**
- `POST /api/v1/files/upload` - Upload document files
- `GET /api/v1/files/{file_id}` - Download/view document

---

### 4. System Health & Monitoring
**Description**: System health checks and monitoring endpoints for operational visibility.

**API:**
- `GET /api/v1/health` - Basic health check


