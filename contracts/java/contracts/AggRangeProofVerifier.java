package contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tuples.generated.Tuple2;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.ContractGasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 3.6.0.
 */
public class AggRangeProofVerifier extends Contract {
    private static final String BINARY = "0x60806040523480156200001157600080fd5b5060405160208062002820833981018060405260208110156200003357600080fd5b8101908080519060200190929190505050806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663da8972246040518163ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040160206040518083038186803b1580156200010857600080fd5b505afa1580156200011d573d6000803e3d6000fd5b505050506040513d60208110156200013457600080fd5b81019080805190602001909291905050506020141515620001bd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f6269732073697a65206e6f7420657175616c000000000000000000000000000081525060200191505060405180910390fd5b620001c76200078c565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166304c09ce96040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401604080518083038186803b1580156200024a57600080fd5b505afa1580156200025f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525060408110156200028557600080fd5b81019080919050509050620002996200078c565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166382529fdb6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401604080518083038186803b1580156200031c57600080fd5b505afa15801562000331573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525060408110156200035757600080fd5b810190809190505090506200036b6200078c565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166324d6147d6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401604080518083038186803b158015620003ee57600080fd5b505afa15801562000403573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525060408110156200042957600080fd5b810190809190505090506200043d620007ae565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634db118756040518163ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016110006040518083038186803b158015620004c257600080fd5b505afa158015620004d7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250611000811015620004fe57600080fd5b8101908091905050905062000512620007ae565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663da2b99ce6040518163ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016110006040518083038186803b1580156200059757600080fd5b505afa158015620005ac573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250611000811015620005d357600080fd5b81019080919050509050846000600281101515620005ed57fe5b60200201516001600001819055508460016002811015156200060b57fe5b602002015160018001819055508360006002811015156200062857fe5b60200201516003600001819055508360016002811015156200064657fe5b60200201516003600101819055508260006002811015156200066457fe5b60200201516005600001819055508260016002811015156200068257fe5b602002015160056001018190555060008090505b60026020028110156200077f578281600202608081101515620006b557fe5b6020020151600782604081101515620006ca57fe5b60020201600001819055508260018260020201608081101515620006ea57fe5b6020020151600782604081101515620006ff57fe5b600202016001018190555081816002026080811015156200071c57fe5b60200201516087826040811015156200073157fe5b600202016000018190555081600182600202016080811015156200075157fe5b60200201516087826040811015156200076657fe5b6002020160010181905550808060010191505062000696565b50505050505050620007d2565b6040805190810160405280600290602082028038833980820191505090505090565b61100060405190810160405280608090602082028038833980820191505090505090565b61203e80620007e26000396000f3fe608060405260043610610083576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680631213e20114610088578063654474ee146100ba5780636d3f6b2714610110578063b7479f5f14610261578063b8c9d365146102b7578063cff0ab96146102e9578063e2179b8e14610340575b600080fd5b34801561009457600080fd5b5061009d610372565b604051808381526020018281526020019250505060405180910390f35b3480156100c657600080fd5b506100f3600480360360208110156100dd57600080fd5b8101908080359060200190929190505050610384565b604051808381526020018281526020019250505060405180910390f35b34801561011c57600080fd5b50610247600480360361052081101561013457600080fd5b81019080806101800190600c806020026040519081016040528092919082600c60200280828437600081840152601f19601f82011690508083019250505050505091929192908060a001906005806020026040519081016040528092919082600560200280828437600081840152601f19601f8201169050808301925050505050509192919290806101800190600c806020026040519081016040528092919082600c60200280828437600081840152601f19601f8201169050808301925050505050509192919290806101800190600c806020026040519081016040528092919082600c60200280828437600081840152601f19601f82011690508083019250505050505091929192905050506103ad565b604051808215151515815260200191505060405180910390f35b34801561026d57600080fd5b5061029a6004803603602081101561028457600080fd5b8101908080359060200190929190505050610638565b604051808381526020018281526020019250505060405180910390f35b3480156102c357600080fd5b506102cc610661565b604051808381526020018281526020019250505060405180910390f35b3480156102f557600080fd5b506102fe610673565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561034c57600080fd5b50610355610698565b604051808381526020018281526020019250505060405180910390f35b60058060000154908060010154905082565b60878160408110151561039357fe5b600202016000915090508060000154908060010154905082565b60006103b7611c54565b6040805190810160405280876000600c811015156103d157fe5b60200201518152602001876001600c811015156103ea57fe5b602002015181525081600001819052506040805190810160405280876002600c8110151561041457fe5b60200201518152602001876003600c8110151561042d57fe5b602002015181525081602001819052506040805190810160405280876004600c8110151561045757fe5b60200201518152602001876005600c8110151561047057fe5b602002015181525081604001819052506040805190810160405280876006600c8110151561049a57fe5b60200201518152602001876007600c811015156104b357fe5b602002015181525081606001819052508460006005811015156104d257fe5b60200201518160800181815250508460016005811015156104ef57fe5b60200201518160a001818152505084600260058110151561050c57fe5b60200201518160c0018181525050610522611cb8565b6040805190810160405280886008600c8110151561053c57fe5b60200201518152602001886009600c8110151561055557fe5b602002015181525081600060028110151561056c57fe5b6020020181905250604080519081016040528088600a600c8110151561058e57fe5b6020020151815260200188600b600c811015156105a757fe5b60200201518152508160016002811015156105be57fe5b6020020181905250848260e0015160000181905250838260e00151602001819052508560036005811015156105ef57fe5b60200201518260e00151604001818152505085600460058110151561061057fe5b60200201518260e00151606001818152505061062c81836106aa565b92505050949350505050565b60078160408110151561064757fe5b600202016000915090508060000154908060010154905082565b60038060000154908060010154905082565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60018060000154908060010154905082565b60006106b4611ce6565b6106bc6113d2565b90506106c6611ce6565b6106ce61146d565b90506106d8611d15565b610704856000015160000151866000015160200151876020015160000151886020015160200151611508565b8160a001818152505061071a8160a0015161155f565b816101000181815250506107318160a0015161159b565b8160c0018190525061074e6107498260a00151611643565b61159b565b8160e00181905250610764816101000151611713565b8161016001818152505061079361078e82610100015183610100015161173f90919063ffffffff16565b61177b565b816101200181815250506107c26107bd82610100015183610120015161173f90919063ffffffff16565b61177b565b816101400181815250506107d660026117b5565b81610180018190525061080b856040015160000151866040015160200151876060015160000151886060015160200151611508565b81606001818152505061082f8160600151826060015161173f90919063ffffffff16565b8160800181815250506108f76108568260800151876060015161185a90919063ffffffff16565b6108e96108748460600151896040015161185a90919063ffffffff16565b6108db6108a28661014001518c600160028110151561088f57fe5b602002015161185a90919063ffffffff16565b6108cd8761012001518d60006002811015156108ba57fe5b602002015161185a90919063ffffffff16565b61190790919063ffffffff16565b61190790919063ffffffff16565b61190790919063ffffffff16565b816101c001819052506109716109286109148361018001516119a9565b83610140015161173f90919063ffffffff16565b6109636109388460c00151611a10565b610955856101200151866101000151611a7a90919063ffffffff16565b61173f90919063ffffffff16565b611a7a90919063ffffffff16565b816102000181815250506109ce6109ba61098f8361018001516119a9565b6109ac84610100015185610140015161173f90919063ffffffff16565b61173f90919063ffffffff16565b826102000151611a7a90919063ffffffff16565b81610200018181525050610a6e610a128660a0015160036040805190810160405290816000820154815260200160018201548152505061185a90919063ffffffff16565b610a60610a318461020001518960800151611a7a90919063ffffffff16565b60016040805190810160405290816000820154815260200160018201548152505061185a90919063ffffffff16565b61190790919063ffffffff16565b816101e00181905250806101e0015160000151816101c0015160000151141580610aaa5750806101e0015160200151816101c001516020015114155b15610abb57600093505050506113cc565b610ac3611e38565b610af4610ae18360600151886020015161185a90919063ffffffff16565b876000015161190790919063ffffffff16565b905060008090505b6001600501811015610d68576000610b908860e001516000015183600202600c81101515610b2657fe5b60200201518960e001516000015160018560020201600c81101515610b4757fe5b60200201518a60e001516020015185600202600c81101515610b6557fe5b60200201518b60e001516020015160018760020201600c81101515610b8657fe5b6020020151611508565b90508084610220015183600681101515610ba657fe5b602002018181525050610bca610bc5828361173f90919063ffffffff16565b61177b565b84610260015183600681101515610bdd57fe5b602002018181525050610c0684610260015183600681101515610bfc57fe5b6020020151611643565b84610280015183600681101515610c1957fe5b60200201818152505060408051908101604052808960e001516000015184600202600c81101515610c4657fe5b602002015181526020018960e001516000015160018560020201600c81101515610c6c57fe5b60200201518152508460400181905250610cbc610cad85610260015184600681101515610c9557fe5b6020020151866040015161185a90919063ffffffff16565b8461190790919063ffffffff16565b925060408051908101604052808960e001516020015184600202600c81101515610ce257fe5b602002015181526020018960e001516020015160018560020201600c81101515610d0857fe5b60200201518152508460400181905250610d58610d4985610280015184600681101515610d3157fe5b6020020151866040015161185a90919063ffffffff16565b8461190790919063ffffffff16565b9250508080600101915050610afc565b5060008090505b6002602002811015611252576000811415610eb05760008090505b6001600501811015610e3a57600084610220015182600681101515610dab57fe5b602002015190506000821415610ddd5780856102c0015184604081101515610dcf57fe5b602002018181525050610e2c565b610e0f610e0a82876102c0015186604081101515610df757fe5b602002015161173f90919063ffffffff16565b61177b565b856102c0015184604081101515610e2257fe5b6020020181815250505b508080600101915050610d8a565b50826102c0015181604081101515610e4e57fe5b602002015183610300015182604081101515610e6657fe5b602002018181525050610e8f836102c0015182604081101515610e8557fe5b6020020151611643565b836102c0015182604081101515610ea257fe5b602002018181525050610fb4565b6000610ec0826001600501611abf565b9050610f1e610f1985610260015183600160050103600681101515610ee157fe5b6020020151866102c00151610ef860018603611b0f565b8603604081101515610f0657fe5b602002015161173f90919063ffffffff16565b61177b565b846102c0015183604081101515610f3157fe5b602002018181525050610f96610f9185610280015183600160050103600681101515610f5957fe5b6020020151866103000151610f7060018603611b0f565b8603604081101515610f7e57fe5b602002015161173f90919063ffffffff16565b61177b565b84610300015183604081101515610fa957fe5b602002018181525050505b826102c0015181604081101515610fc757fe5b6020020151836102a0015182604081101515610fdf57fe5b60200201818152505082610300015181604081101515610ffb57fe5b6020020151836102e001518260408110151561101357fe5b6020020181815250506110658361010001516110578960e0015160400151866102a001518560408110151561104457fe5b602002015161173f90919063ffffffff16565b611b5790919063ffffffff16565b836102a001518260408110151561107857fe5b6020020181815250506110b38760e0015160600151846102e00151836040811015156110a057fe5b602002015161173f90919063ffffffff16565b836102e00151826040811015156110c657fe5b60200201818152505060208110156111515761113061110a846101800151836020811015156110f157fe5b602002015185610120015161173f90919063ffffffff16565b846102e001518360408110151561111d57fe5b6020020151611a7a90919063ffffffff16565b836102e001518260408110151561114357fe5b6020020181815250506111d2565b6111b561118f84610180015160208481151561116957fe5b0660208110151561117657fe5b602002015185610140015161173f90919063ffffffff16565b846102e00151836040811015156111a257fe5b6020020151611a7a90919063ffffffff16565b836102e00151826040811015156111c857fe5b6020020181815250505b61122983610100015161121b8560e00151846040811015156111f057fe5b6020020151866102e001518560408110151561120857fe5b602002015161173f90919063ffffffff16565b611a7a90919063ffffffff16565b836102e001518260408110151561123c57fe5b6020020181815250508080600101915050610d6f565b50600061128c876080015160405160200180828152602001915050604051602081830303815290604052805190602001206001900461177b565b90506113936112c88860c0015160036040805190810160405290816000820154815260200160018201548152505061185a90919063ffffffff16565b61138561134861131961130a8c608001516112fc8e60e00151606001518f60e001516040015161173f90919063ffffffff16565b611a7a90919063ffffffff16565b8661173f90919063ffffffff16565b60056040805190810160405290816000820154815260200160018201548152505061185a90919063ffffffff16565b61137761135a89896102e00151611b93565b6113698b8a6102a00151611b93565b61190790919063ffffffff16565b61190790919063ffffffff16565b61190790919063ffffffff16565b836101e001819052508160000151836101e00151600001511480156113c457508160200151836101e0015160200151145b955050505050505b92915050565b6113da611ce6565b6113e2611ce6565b60008090505b60026020028110156114655760078160408110151561140357fe5b6002020160000154828260408110151561141957fe5b6020020151600001818152505060078160408110151561143557fe5b6002020160010154828260408110151561144b57fe5b6020020151602001818152505080806001019150506113e8565b508091505090565b611475611ce6565b61147d611ce6565b60008090505b60026020028110156115005760878160408110151561149e57fe5b600202016000015482826040811015156114b457fe5b602002015160000181815250506087816040811015156114d057fe5b600202016001015482826040811015156114e657fe5b602002015160200181815250508080600101915050611483565b508091505090565b60006115558585858560405160200180858152602001848152602001838152602001828152602001945050505050604051602081830303815290604052805190602001206001900461177b565b9050949350505050565b60006115948260405160200180828152602001915050604051602081830303815290604052805190602001206001900461177b565b9050919050565b6115a3611e52565b60018160006040811015156115b457fe5b602002018181525050818160016040811015156115cd57fe5b6020020181815250506000600290505b600260200281101561163d5761161961161484846001850360408110151561160157fe5b602002015161173f90919063ffffffff16565b61177b565b828260408110151561162757fe5b60200201818152505080806001019150506115dd565b50919050565b6000807f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000019050600083905060008114156116825760009250505061170e565b8181111561169957818181151561169557fe5b0690505b600080600190506000849050600084905060005b6000821415156116e65781838115156116c257fe5b049050838482028603838484028603809550819650829750839850505050506116ad565b60008512156117035784600003870397505050505050505061170e565b849750505050505050505b919050565b6000817f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001039050919050565b6000807f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000190508080151561176f57fe5b83850991505092915050565b6000807f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001905080838115156117ac57fe5b06915050919050565b6117bd611e76565b60018160006020811015156117ce57fe5b602002018181525050818160016020811015156117e757fe5b6020020181815250506000600290505b60208110156118545761183061182b84846001850360208110151561181857fe5b602002015161173f90919063ffffffff16565b61177b565b828260208110151561183e57fe5b60200201818152505080806001019150506117f7565b50919050565b611862611e38565b600182141561187357829050611901565b600282141561188d576118868384611907565b9050611901565b611895611e9a565b83600001518160006003811015156118a957fe5b60200201818152505083602001518160016003811015156118c657fe5b602002018181525050828160026003811015156118df57fe5b6020020181815250506040826060836007600019fa15156118ff57600080fd5b505b92915050565b61190f611e38565b611917611ebd565b836000015181600060048110151561192b57fe5b602002018181525050836020015181600160048110151561194857fe5b602002018181525050826000015181600260048110151561196557fe5b602002018181525050826020015181600360048110151561198257fe5b6020020181815250506040826080836006600019fa15156119a257600080fd5b5092915050565b6000808260006020811015156119bb57fe5b602002015190506000600190505b6020811015611a06576119f784826020811015156119e357fe5b602002015183611b5790919063ffffffff16565b915080806001019150506119c9565b5080915050919050565b600080826000604081101515611a2257fe5b602002015190506000600190505b6002602002811015611a7057611a618482604081101515611a4d57fe5b602002015183611b5790919063ffffffff16565b91508080600101915050611a30565b5080915050919050565b6000807f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001905082841015611ab2578383820301611ab6565b8284035b91505092915050565b6000808260019060020a02905060005b8185108015611ade5750600082115b15611aff57600182908060020a820491505091508080600101915050611acf565b8084600101039250505092915050565b6000808290506000811415611b28576001915050611b52565b6000600190505b6000821115611b4c57600281029050818060019003925050611b2f565b80925050505b919050565b6000807f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001905080801515611b8757fe5b83850891505092915050565b611b9b611e38565b611ba3611e38565b611bdc836000604081101515611bb557fe5b6020020151856000604081101515611bc957fe5b602002015161185a90919063ffffffff16565b90506000600190505b6002602002811015611c4957611c3a611c2b8583604081101515611c0557fe5b60200201518784604081101515611c1857fe5b602002015161185a90919063ffffffff16565b8361190790919063ffffffff16565b91508080600101915050611be5565b508091505092915050565b6104a060405190810160405280611c69611ee0565b8152602001611c76611ee0565b8152602001611c83611ee0565b8152602001611c90611ee0565b8152602001600081526020016000815260200160008152602001611cb2611efa565b81525090565b6080604051908101604052806002905b611cd0611ee0565b815260200190600190039081611cc85790505090565b611000604051908101604052806040905b611cff611ee0565b815260200190600190039081611cf75790505090565b6160c060405190810160405280611d2a611f30565b8152602001611d37611f30565b8152602001611d44611ee0565b8152602001600081526020016000815260200160008152602001611d66611f54565b8152602001611d73611f54565b815260200160008152602001600081526020016000815260200160008152602001611d9c611f78565b8152602001611da9611f9c565b8152602001611db6611ee0565b8152602001611dc3611ee0565b815260200160008152602001611dd7611fcb565b8152602001611de4611fcb565b8152602001611df1611fcb565b8152602001611dfe611fcb565b8152602001611e0b611f54565b8152602001611e18611f54565b8152602001611e25611f54565b8152602001611e32611f54565b81525090565b604080519081016040528060008152602001600081525090565b61080060405190810160405280604090602082028038833980820191505090505090565b61040060405190810160405280602090602082028038833980820191505090505090565b606060405190810160405280600390602082028038833980820191505090505090565b608060405190810160405280600490602082028038833980820191505090505090565b604080519081016040528060008152602001600081525090565b61034060405190810160405280611f0f611fee565b8152602001611f1c611fee565b815260200160008152602001600081525090565b61100060405190810160405280608090602082028038833980820191505090505090565b61080060405190810160405280604090602082028038833980820191505090505090565b61040060405190810160405280602090602082028038833980820191505090505090565b610800604051908101604052806020905b611fb5611ee0565b815260200190600190039081611fad5790505090565b60c060405190810160405280600690602082028038833980820191505090505090565b61018060405190810160405280600c9060208202803883398082019150509050509056fea165627a7a72305820d59a53a226be4e7f01de08ee1de4db992b69581ca3cc759bfe956d73bfee26d70029";

    public static final String FUNC_UBASE = "uBase";

    public static final String FUNC_HVBASE = "hvBase";

    public static final String FUNC_GVBASE = "gvBase";

    public static final String FUNC_H = "h";

    public static final String FUNC_PARAMS = "params";

    public static final String FUNC_G = "g";

    public static final String FUNC_AGGVERIFYRANGEPROOF = "aggVerifyRangeProof";

    @Deprecated
    protected AggRangeProofVerifier(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected AggRangeProofVerifier(String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected AggRangeProofVerifier(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected AggRangeProofVerifier(String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> uBase() {
        final Function function = new Function(FUNC_UBASE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<BigInteger, BigInteger>>(
                new Callable<Tuple2<BigInteger, BigInteger>>() {
                    @Override
                    public Tuple2<BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> hvBase(BigInteger param0) {
        final Function function = new Function(FUNC_HVBASE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<BigInteger, BigInteger>>(
                new Callable<Tuple2<BigInteger, BigInteger>>() {
                    @Override
                    public Tuple2<BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> gvBase(BigInteger param0) {
        final Function function = new Function(FUNC_GVBASE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<BigInteger, BigInteger>>(
                new Callable<Tuple2<BigInteger, BigInteger>>() {
                    @Override
                    public Tuple2<BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> h() {
        final Function function = new Function(FUNC_H, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<BigInteger, BigInteger>>(
                new Callable<Tuple2<BigInteger, BigInteger>>() {
                    @Override
                    public Tuple2<BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<String> params() {
        final Function function = new Function(FUNC_PARAMS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> g() {
        final Function function = new Function(FUNC_G, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<BigInteger, BigInteger>>(
                new Callable<Tuple2<BigInteger, BigInteger>>() {
                    @Override
                    public Tuple2<BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public static RemoteCall<AggRangeProofVerifier> deploy(Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider, String params_) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(params_)));
        return deployRemoteCall(AggRangeProofVerifier.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor);
    }

    public static RemoteCall<AggRangeProofVerifier> deploy(Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider, String params_) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(params_)));
        return deployRemoteCall(AggRangeProofVerifier.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<AggRangeProofVerifier> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, String params_) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(params_)));
        return deployRemoteCall(AggRangeProofVerifier.class, web3j, credentials, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<AggRangeProofVerifier> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, String params_) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(params_)));
        return deployRemoteCall(AggRangeProofVerifier.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    public RemoteCall<Boolean> aggVerifyRangeProof(List<BigInteger> points, List<BigInteger> scalar, List<BigInteger> l, List<BigInteger> r) {
        final Function function = new Function(FUNC_AGGVERIFYRANGEPROOF, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.StaticArray12<org.web3j.abi.datatypes.generated.Uint256>(
                        org.web3j.abi.Utils.typeMap(points, org.web3j.abi.datatypes.generated.Uint256.class)), 
                new org.web3j.abi.datatypes.generated.StaticArray5<org.web3j.abi.datatypes.generated.Uint256>(
                        org.web3j.abi.Utils.typeMap(scalar, org.web3j.abi.datatypes.generated.Uint256.class)), 
                new org.web3j.abi.datatypes.generated.StaticArray12<org.web3j.abi.datatypes.generated.Uint256>(
                        org.web3j.abi.Utils.typeMap(l, org.web3j.abi.datatypes.generated.Uint256.class)), 
                new org.web3j.abi.datatypes.generated.StaticArray12<org.web3j.abi.datatypes.generated.Uint256>(
                        org.web3j.abi.Utils.typeMap(r, org.web3j.abi.datatypes.generated.Uint256.class))), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    @Deprecated
    public static AggRangeProofVerifier load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new AggRangeProofVerifier(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static AggRangeProofVerifier load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new AggRangeProofVerifier(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static AggRangeProofVerifier load(String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
        return new AggRangeProofVerifier(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static AggRangeProofVerifier load(String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        return new AggRangeProofVerifier(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
