/**
 * campaign contract is used for rnodes to claim campaign.
 * the general campaign process is as follows:
 * 1. nodes submit parameters(numOfCampaign) and proofs, including memory and cpu proofs.
 * 2. each proof consists of an address, a nonce and a block number(hash);
 * 3. 'address' prevents nodes from stealing others' proof, 'block hash' avoids calculation in advance;
 * 4. campaign contract checks if all requirements are satisfied:
 *    rnode, admission and parameters
 * 5. if pass, nodes become candidate.
**/


pragma solidity ^0.4.24;

import "./lib/safeMath.sol";
import "./lib/set.sol";

// there are two interfaces to interact with admission and rnode contracts
// rnode and admission contracts must be deployed before this campaign contract, because
// addresses of rnode and admission contracts are needed to initiate interfaces during construction of campaign contract
contract AdmissionInterface {
    function verify(
        uint64 _cpuNonce,
        uint _cpuBlockNumber,
        uint64 _memoryNonce,
        uint _memoryBlockNumber,
        address _sender
    )
    public
    view
    returns (bool);
}

contract RnodeInterface{
    function isRnode(address _addr)public view returns (bool);
}


// TODO this file exposes security holes.

contract Campaign {

    using Set for Set.Data;
    using SafeMath for uint256;

    address owner; // owner has permission to set parameters
    uint public termIdx = 0; // current term
    uint public viewLen = 3; // view length: the number of blocks it can propose within a term
    uint public termLen =12; // term length: the number of proposers in a certain term
    // 'round' is same with 'term'.
    uint public numPerRound = termLen * viewLen; // total number of blocks produced in a term.
    // a node must choose the number of terms when it claims campaign
    uint public minNoc = 1; // minimal number of terms
    uint public maxNoc = 10; //maximum number of terms

    uint public withdrawTermIdx = 0; // withdraw deposit after each round.
    bool withdrawFlag = true;

    // a new type for a single candidate
    struct CandidateInfo {
        uint numOfCampaign; // rest terms that the candidate will claim campaign
        uint startTermIdx;
        uint stopTermIdx;
    }

    mapping(address => CandidateInfo) candidates; // store a 'CandidateInfo' struct for each possible address
    mapping(uint => Set.Data) campaignSnapshots; // store all candidates for each term

    // declare admission and rnode interfaces
    AdmissionInterface admission;
    RnodeInterface    rnode;

    modifier onlyOwner() {require(msg.sender == owner);_;}

    event ClaimCampaign(address candidate, uint startTermIdx, uint stopTermIdx);
    event QuitCampaign(address candidate, uint payback);
    event ViewChange();

    // admission and rnode interfaces will be initiated during creation of campaign contract
    constructor(address _admissionAddr, address _rnodeAddr) public {
        owner = msg.sender;
        admission = AdmissionInterface(_admissionAddr);
        rnode = RnodeInterface(_rnodeAddr);
        withdrawTermIdx = (block.number - 1).div(numPerRound);
    }

    function() payable public { }

    // get all candidates of given term index
    function candidatesOf(uint _termIdx) public view returns (address[]){
        return campaignSnapshots[_termIdx].values;
    }

    // get info of given candidate
    function candidateInfoOf(address _candidate)
    public
    view
    returns (uint, uint, uint)
    {
        return (
        candidates[_candidate].numOfCampaign,
        candidates[_candidate].startTermIdx,
        candidates[_candidate].stopTermIdx
        );
    }

    // admission interface can be set afterwards by owner
    function setAdmissionAddr(address _addr) public onlyOwner(){
        admission = AdmissionInterface(_addr);
    }

    // rnode interface can be set afterwards by owner
    function setRnodeInterface(address _addr) public onlyOwner(){
        rnode = RnodeInterface(_addr);
    }

    // owner can set these parameters
    function updateMinNoc(uint _minNoc) public onlyOwner(){
        minNoc = _minNoc;
    }

    function updateMaxNoc(uint _maxNoc) public onlyOwner(){
        maxNoc = _maxNoc;
    }

    function updateTermLen(uint _termLen) public onlyOwner(){
        termLen = _termLen;
        numPerRound = SafeMath.mul(termLen, viewLen);
    }

    function updateViewLen(uint _viewLen) public onlyOwner(){
        viewLen = _viewLen;
        numPerRound = SafeMath.mul(termLen, viewLen);
    }

    /**
     * Submits required information to participate the campaign for membership of the committee.
     *
     * Each call may tried to update termIdx once.
     *
     * Claiming a candidate has three conditions:
     * 1. pay some specified cpc token.
     * 2. pass the cpu&memory capacity proof.
     * 3. rpt score reaches the threshold (be computed somewhere else).
     */

    // claim campaign will verify these parameters, if pass, the node will become candidate
    function claimCampaign(
        uint _numOfCampaign,  // number of terms that the node want to claim campaign

        // admission parameters
        uint64 _cpuNonce,
        uint _cpuBlockNumber,
        uint64 _memoryNonce,
        uint _memoryBlockNumber
    )
    public
    payable
    {
        // initiate withdrawTermIdx during first call,
        // in case that termIdx too large while withdrawTermIdx too low,
        // resulting in large 'for' loop and gas not enough.
        if(withdrawFlag) {
            withdrawTermIdx = (block.number - 1).div(numPerRound) - 10;
            withdrawFlag = false;
        }
        // only rnode can become candidate
        require(rnode.isRnode(msg.sender)==true, "not RNode by rnode");

        // verify the sender's cpu&memory ability.
        require(admission.verify(_cpuNonce, _cpuBlockNumber, _memoryNonce, _memoryBlockNumber, msg.sender), "cpu or memory not passed.");
        require((_numOfCampaign >= minNoc && _numOfCampaign <= maxNoc), "num of campaign out of range.");

        updateCandidateStatus(); // update status firstly,  then check

        address candidate = msg.sender;

        // if nodes have not ended their terms, they can not claim again
        require(
            candidates[candidate].numOfCampaign == 0,
            "please waite until your last round ended and try again."
        );

        // set candidate's numOfCampaign according to arguments, and set start and end termIdx respectively
        candidates[candidate].numOfCampaign = _numOfCampaign;
        candidates[candidate].startTermIdx = termIdx.add(1);

        //[start, stop)
        candidates[candidate].stopTermIdx = candidates[candidate].startTermIdx.add(_numOfCampaign);

        // add candidate to campaignSnapshots.
        for(uint i = candidates[candidate].startTermIdx; i < candidates[candidate].stopTermIdx; i++) {
            campaignSnapshots[i].insert(candidate);
        }
        emit ClaimCampaign(candidate, candidates[candidate].startTermIdx, candidates[candidate].stopTermIdx);
    }

    /**
     * The function will be called when a node claims to campaign for proposer election to update candidates status.
     *
     */

    // update candidate status, i.e. subtract 1 from numOfCampaign when a new term begin
    // withdrawTermIdx record the start term that need to update
    function updateCandidateStatus() public payable {
        // get current term, update termIdx
        updateTermIdx();

        if (withdrawTermIdx >= termIdx) {
            return;
        }

        uint size;
        // withdrawTermIdx is the last term where candidates claim campaign
        for(; withdrawTermIdx <= termIdx; withdrawTermIdx++) {
            // avoid recalculate the size for circulation times.
            size = campaignSnapshots[withdrawTermIdx].values.length;
            // go through all candidates in term withdrawTermIdx, and update their numOfCampaign
            for(uint i = 0; i < size; i++) {
                address candidate = campaignSnapshots[withdrawTermIdx].values[i];

                if (candidates[candidate].numOfCampaign == 0) {
                    continue;
                }

                candidates[candidate].numOfCampaign = SafeMath.sub(candidates[candidate].numOfCampaign, 1);

                // if candidate's tenure is all over, all status return to zero.
                if (candidates[candidate].numOfCampaign == 0) {
                    candidates[candidate].startTermIdx = 0;
                    candidates[candidate].stopTermIdx = 0;
                }
            }
        }
    }

    /** update termIdx called by function ClaimCampaign. */
    function updateTermIdx() internal{
        uint blockNumber = block.number;
        if (blockNumber == 0) {
            termIdx = 0;
            return;
        }

        // calculate current term
        termIdx = (blockNumber - 1).div(numPerRound);
    }

}
