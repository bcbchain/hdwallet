# HD钱包使用手册

> HD钱包是一个独立的命令行程序，提供了对rpc接口的本地调用封装。便于日常性的诊断、集成，不需要再去构造POST请求。
>
> `HD Wallets`的全称是 `Hierachical Deterministic Wallets`， 对应中文是 分层确定性钱包。

## 1. 使用方法

   首先修改`hdWallet`配置文件`.config/hdWallet.yaml`，如果在测试链上使用只需要修改`chainID`为测试链`chainID`，`nodeAddrSlice`改为测试链的接入的URL就可以了。（URL不是域名的需要改为<http://ip:46657>, `rpc`服务如果使用的是`https`需要将`useHttps`改为`true`）

   然后使用hdWallet_rpc 命令行工具生成助记词，然后导出助记词，助记词必须备份好，如果丢失账户将无法找回，后果自负。

​   最后使用hdWallet命令行工具生成钱包，然后进行相应的操作。

## 2. 配置文件说明

.config/hdWallet.yaml 为配置文件，如下

```yaml
#区块链标识
chainID: "bcb"

#区块链版本
chainVersion: 1

# 本地服务监听地址:端口
serverAddr: "tcp://0.0.0.0:37657"

# 是否使用https
useHttps: false

# 外部证书路径
outCerPath: "./.config/Xwallet.web.rpc"

#日志配置信息
loggerScreen: true
loggerFile: true
loggerLevel: "debug"

# 账户数据库存位置
keyStorePath: "./.keystore"

#指定创智区块链节点的接入URL，热钱包需要此参数，冷钱包可以忽略此参数；
nodeAddrSlice:
    - "https://earth.bcbchain.io"
    - "https://mars.bcbchain.io"
    - "https://mercury.bcbchain.io"
    - "https://jupiter.bcbchain.io"
    - "https://venus.bcbchain.io"
    - "https://moon.bcbchain.io"
    - "https://sirius.bcbchain.io"
    - "https://vaga.bcbchain.io"
    - "https://altair.bcbchain.io"
```

| **选项**      | **类型** | **注释**                                                     |
| ------------- | -------- | ------------------------------------------------------------ |
| chainID       | String   | 链id，默认bcb。                                              |
| chainVersion  | Int      | 设置为1时，发送v1版本交易。<br />设置为2时，发送v2版本交易。 |
| serverAddr    | String   | 启动rpc服务时本地服务监听地址:端口                           |
| useHttps      | String   | 是否使用https，默认为false                                   |
| outCerPath    | String   | 外部证书路径                                                 |
| loggerScreen  | Bool     | 日志是否打印在屏幕上                                         |
| loggerFile    | Bool     | 日志是否打印在文件里                                         |
| loggerLevel   | String   | 日志等级                                                     |
| keyStorePath  | String   | 账户数据库存位置                                             |
| nodeAddrSlice | String   | 指定创智区块链节点的接入URL，<br />热钱包需要此参数，冷钱包可以忽略此参数； |

## 3. 命令详解

### 3.1 hdWallet_rpc命令

hdWallet_rpc 命令运行格式如下：

```bash
Usage:
  hdWallet_rpc [command]

Available Commands:
  changePassword Change password
  createMnemonic Create mnemonic
  exportMnemonic Export mnemonic
  help           Help about any command
  importMnemonic Import mnemonic

Flags:
  -h, --help   help for hdWallet_rpc

Use "hdWallet_rpc [command] --help" for more information about a command.
```

#### 3.1.1 createMnemonic

> 注：危险！：在使用HD钱包生成助记词之后必须进行助记词备份，否则钱包丢失或忘记密码后无法恢复钱包。
> 注：因为一个HD钱包可生成足够多私钥的属性，以及方便使用和管理，本版本只支持创建一次HD钱包。

- **command**

  ```bash
  hdWallet_rpc createMnemonic
  ```

- **Output SUCCESS Example**

  ```json
  {
    "mnemonic": "craft toilet twist safe violin catch similar add friend lion fabric crisp",
    "password": "9FPqXE82XLHwmoqK2TNXV9Lg45fhXcUREESzoK7JCcnX"
  }
  
  ```

- **Output SUCCESS Result**

  | **语法** | **类型** | **注释**                                                      |
  | -------- | ------ | ------------------------------------------------------------ |
  | mnemonic | String | 12个助记词，以单个空格隔开。                                 |
  | password | String | 随机生成的新助记词密码                                       |

#### 3.1.2 exportMnemonic

- **command**

  ```bash
  hdWallet_rpc exportMnemonic --password 9FPqXE82XLHwmoqK2TNXV9Lg45fhXcUREESzoK7JCcnX
  ```

- **Input Parameters**

  | **语法** | **类型** | **注释**                                                     |
  | -------- | ------ | ------------------------------------------------------------ |
  | password | String | 助记词密码。                                                 |

- **Output SUCCESS Example**

  ```json
  {
    "mnemonic": "craft toilet twist safe violin catch similar add friend lion fabric crisp"
  }
  
  ```

- **Output SUCCESS Result**

  | **语法** | **类型** | **注释**                                                      |
  | -------- | ------ | ------------------------------------------------------------ |
  | mnemonic | String | 12个助记词，以单个空格隔开。                                 |

#### 3.1.3 importMnemonic

> 注：密码遗失后，只能重新导入助记词生成新的密码。

- **command**

  ```bash
  hdWallet_rpc importMnemonic --mnemonic “craft toilet ... crisp”
  ```

- **Input Parameters**

  | **语法** | **类型** | **注释**                                                      |
  | -------- | ------ | ------------------------------------------------------------ |
  | mnemonic | String | 12个助记词，以单个空格隔开。                                 |

- **Output SUCCESS Example**

  ```json
  {
   "password":2oHQjWWKkGfUvWpeXV8qwgCTzrDh9fJ8hJV2L6gP6JeU
  }
  
  ```

- **Output SUCCESS Result**

  | **语法** | **类型** | **注释**                                                      |
  | -------- | ------ | ------------------------------------------------------------ |
  | password | String | 随机生成新的助记词密码。                                     |

#### 3.1.4 changePassword

- **command**

  ```bash
  hdWallet_rpc changePassword --password 2oHQjWWKkGfUvWpeXV8qwgCTzrDh9fJ8hJV2L6gP6JeU
  ```

- **Input Parameters**

  | **语法** | **类型** | **注释**                                                      |
  | -------- | ------ | ------------------------------------------------------------ |
  | password | String | 助记词密码。                                                  |

- **Output SUCCESS Example**

  ```json
  {
    "password":6eaN3vQjEW1ncCu4fYc9AvFFxFE4f4sY5HKw31X35zZ4
  }
  
  ```

- **Output SUCCESS Result**

  | **语法** | **类型** | **注释** |
  | -------- | ------ | ------------------------------------------------------------ |
  | password |  String  | 随机生成的新助记词密码。                                     |

### 3.2 hdWallet

hdWalle 命令运行格式如下：

```bash
Usage:
  hdWallet [command]

Available Commands:
  allBalance      get balance of all tokens for specific address
  balance         get balance of BCB token for specific address
  balanceOfToken  get balance of specific token for specific address
  block           get block info with height
  blockHeight     get current block height
  commitTx        commit transaction
  help            Help about any command
  nonce           get the next usable nonce for specific address
  transaction     get transaction info with txHash
  transfer        transfer token
  transferOffline offline pack and sign transfer transaction
  walletCreate    create wallet

Flags:
  -h, --help   help for hdWallet

Use "hdWallet [command] --help" for more information about a command.
```

#### 3.2.1 walletCreate

> 注：使用HD钱包时，钱包是根据path作为参数之一生成的私钥，生成的私钥同name没有关系，外部要做好name和path的对应关系，以免name和私钥对应不上。

- **command**

  ```bash
  hdWallet walletCreate --path "m/44'/60'/0'/0/1" --password aBig62_123 [--url https://...]
  ```

- **Input Parameters**

  | **选项** | **类型** | **注释**                                                      |
  | -------- | ------ | ------------------------------------------------------------ |
  | path     | String  | HD钱包路径:"m/44'/60'/0'/0/#", 其中”#“为数字，范围：[0, 4294967295]。 |
  | password | String | 助记词密码。                                                         |
  | url      | String | 钱包服务地址，可选项，默认调用本地服务。                               |

- **Output FAILED Example**

  ```json
  {
    "code": -32603,
    "message": "Invalid parameters.",
    "data": ""
  }
  ```

  注：所有命令执行错误返回格式相同，后面不再进行说明。

- **Output SUCCESS Example**

  ```json
  {
    "walletAddr": "bcbES5d6kwoX4vMeNLENMee2Mnsf2KL9ZpWo"
  }
  ```

- **Output SUCCESS Result**

  | **语法**   | **类型** | **注释**                                                      |
  | ---------- | ------- | ------------------------------------------------------------ |
  | walletAddr | Address  | 钱包地址。                                                   |

#### 3.2.2 transfer

- **command**

  ```bash
  hdWallet transfer --path "m/44'/60'/0'/0/1" --password BYd.. --smcAddress bcbLVgb... --gasLimit 600 [--note hello] --to bcbLocFJG5Q792eLQXhvNkG417kwiaaoPH5a --value 1500000000 [--url https://...]
  ```

- **Request Parameters**

  | **语法**   | **类型** | **注释**                                                      |
  | ---------- | ------- | ------------------------------------------------------------ |
  | password   | String  | 助记词密码。                                                  |
  | path       | String  | HD钱包路径:"m/44'/60'/0'/0/#", 其中”#“为数字。                |
  | smcAddress | Address | 交易资产（本币或代币）的token地址。                            |
  | gasLimit   | String  | 交易的燃料限制。                                              |
  | note       | String  | 交易备注（最长256字符），可选项。                              |
  | to         | Address | 接收转账的账户地址。                                          |
  | value      | String  | 转账的资产数量（单位：Cong）。                                 |
  | url        | String  | 钱包服务地址，可选项，默认调用本地服务。                        |

- **Output SUCCESS Example**

  ```json
  {
    "code": 200,
    "log": "Check tx succeed",
    "txHash": "0xA1C960B9D5DB633A6E45B45015A722A2C516B392F93C9BF41F5DAA1197030584",
    "height": 234389
  }
  ```

- **Output SUCCESS Result**

  | **语法**           | **类型**  | **注释**                                                         |
  | ------------------ | --------- | --------------------------------------------------------------- |
  | &nbsp;&nbsp;code   | Int       | 交易执行结果代码，200表示成功。                                   |
  | &nbsp;&nbsp;log    | String    | 对交易执行结果进行的文字描述，当code不等于200时描述具体的错误信息。 |
  | &nbsp;&nbsp;txHash | HexString | 交易的哈希，以 0x 开头。                                         |
  | &nbsp;&nbsp;height | Int64     | 交易在哪个高度的区块被确认。                                      |

#### 3.2.3 transferOffline

- **command**

  ```bash
  hdWallet transferOffline --path "m/44'/60'/0'/0/1" --password BYd... --smcAddress bcbLVgb... --gasLimit 600 [--note hello] --nonce 1500 --to bcbLocF... --value 1500000000 [--url https://...]
  ```

- **Request Parameters**

  | **语法**   | **类型** | **注释**                                                     |
  | ---------- | ------- | ------------------------------------------------------------ |
  | password   | String  | 助记词密码。                                                 |
  | path       | String  | HD钱包路径:"m/44'/60'/0'/0/#", 其中”#“为数字。                |
  | smcAddress | Address | 交易资产（本币或代币）的token地址。                           |
  | gasLimit   | String  | 交易的燃料限制。                                             |
  | note       | String  | 交易备注（最长256字符），可选项。                             |
  | nonce      | String  | 交易计数值。                                                 |
  | to         | Address | 接收转账的账户地址。                                          |
  | value      | String  | 转账的资产数量（单位：Cong）。                                |
  | url        | String  | 钱包服务地址，可选项，默认调用本地服务。                       |

- **Output SUCCESS Example**

  ```json
  {
    "tx": "bcb<tx>.v2.AetboYAmy2TEyUbsR731FTLDLyHE1MVKsSd4v7hS1jFnNkrtmGEVxVmWHR3jVSU
    ffxKgW7aPawnQaVrZ4gwMt6aogUAJjhvnukfPWnxmsybqDgdjgecjsXa94bamPqgPhTTZC9Szb.<1>.YT
    giA1gdDGi2L8iCryAn34dXVYKUEdmBxivyHbK57wKpBcX5KrKyn1vdmZTuKKZ7PotCjcbASbesv61VLE8
    H38TDiopHrs2eHG9z9iEDDyLcN7giLPCgFiLN9LPRiYZgxwpR95echr2bRPbijnKWj"
  }
  ```

- **Output SUCCESS Result**

  | **语法**       | **类型** | **注释**                                                      |
  | -------------- | ------ | ------------------------------------------------------------ |
  | &nbsp;&nbsp;tx | String | 生成的离线交易数据。                                          |

#### 3.2.4 commitTx

- **command**

  ```bash
  hdWallet commitTx --tx "bcb<tx>.v2.AetboY... [--url https://...]"
  ```

- **Input Parameters**

  | **语法** | **类型** | **注释**                                                   |
  | -------- | ------ | ------------------------------------------------------------ |
  | tx       | String | 交易数据。                                                   |
  | url      | String | 钱包服务地址，可选项，默认调用本地服务。                       |

- **Output SUCCESS Example**

  ```json
  {
     "code": 200,
     "log": "Deliver tx succeed",
     "txHash": "0xA1C960B9D5DB633A6E45B45015A722A2C516B392F93C9BF41F5DAA1197030584"
     "height": 234389
  }
  ```

- **Output SUCCESS Result**

  | **语法**           | **类型**   | **注释**                                                             |
  | ------------------ | --------- | -------------------------------------------------------------------- |
  | &nbsp;&nbsp;code   | Int       | 交易校验/背书结果代码，200表示成功。                                   |
  | &nbsp;&nbsp;log    | String    | 对交易校验/背书结果进行的文字描述，当code不等于200时描述具体的错误信息。 |
  | &nbsp;&nbsp;txHash | HexString | 交易的哈希，以 0x 开头。                                              |
  | &nbsp;&nbsp;height | Int64     | 交易在哪个高度的区块被确认。                                           |

#### 3.2.5 blockHeight

- **command**

  ```bash
  hdWallet blockHeight [--url https://...]
  ```

- **Input Parameters**

  | **语法** | **类型** | **注释**                                                      |
  | -------- | ------ | ------------------------------------------------------------ |
  | url      | String | 钱包服务地址，可选项，默认调用本地服务。                        |

- **Output SUCCESS Example**

  ```json
  {
      "lastBlock": 2500
  }
  ```

- **Output SUCCESS Result**

  | **语法**  | **类型** | **注释**                                                    |
  | --------- | ------ | ------------------------------------------------------------ |
  | lastBlock | Int64  | 最新区块高度。                                                |

#### 3.2.6 block

- **command**

  ```bash
  hdWallet block --height 68685 [--url https://...]
  ```

- **Input Parameters**

  | **语法** | **类型** | **注释**                                                      |
  | -------- | ------ | ------------------------------------------------------------ |
  | height   | int64   | 指定区块高度，为0时返回最新高度的区块信息。                  |
  | url      | String  | 钱包服务地址，可选项，默认调用本地服务。                     |

- **Output SUCCESS Example**

  ```json
  {
    "blockHeight": 68685,
    "blockHash": "0x206e6bcfd17e1a64390e61c69468faee5a11faea",
    "parentHash": "0x05d7d5ff4a6d9a83da4c5b6f5913959a667b3b7b",
    "chainID": "devtest",
    "validatorHash": "0xd0d9a92830046e4abe4959be18ad0d8440dd33af",
    "consensusHash": "0xf66ef1df8ba6dac7a1ecce40cc84e54a1cebc6a5",
    "blockTime": "2019-10-28 08:34:33.134876793 +0000 UTC",
    "blockSize": 1615,
    "proposerAddress": "devtest4JjGSRE7QpWHgixhDTodsih1joE9uHN16",
    "txs": [
      {
        "txHash": "0xc8dbcfbe10fc7e935c2a29dd3aeb8ca08042a654332389d157f347a98231338e",
        "txTime": "2019-10-28 08:34:33.134876793 +0000 UTC",
        "code": 200,
        "log": "",
        "blockHash": "0x206e6bcfd17e1a64390e61c69468faee5a11faea",
        "blockHeight": 68685,
        "from": "devtestQ4suFdGVB4AbDnNLJ2xP9RXmqro6XtQXX",
        "nonce": 49,
        "gasLimit": 1000000,
        "fee": 1250000,
        "note": "transfer",
        "messages": [
          {
            "smcAddress": "devtestLL6sMXu8s2hhFRoZH67Q8fig9djogVi3H",
            "smcName": "token-basic",
            "method": "Transfer(types.Address,bn.Number)",
            "to": "devtestP4qDrEyBZegP5whAMmB1yy3c36HVCAWc",
            "value": "1000000000"
          }
        ],
        "transferReceipts": [
          {
            "token": "devtestLL6sMXu8s2hhFRoZH67Q8fig9djogVi3H",
            "from": "devtestQ4suFdGVB4AbDnNLJ2xP9RXmqro6XtQXX",
            "to": "devtestP4qDrEyBZegP5whAMmB1yy3c36HVCAWc",
            "value": 1000000000,
            "note": "transfer"
          }
        ]
      }
    ]
  }
  ```

- **Output SUCCESS Result**

| **语法**              | **类型**     | **注释**                                              |
| --------------------- | ------------ | ----------------------------------------------------- |
| blockHeight           | Int64        | 区块高度。                                            |
| blockHash             | HexString    | 区块哈希值，以 0x 开头。                               |
| parentHash            | HexString    | 父区块哈希值，以 0x 开头。                             |
| chainID               | String       | 链ID。                                                |
| validatorsHash        | HexString    | 验证者列表哈希值，以 0x 开头。                         |
| consensusHash         | HexString    | 共识信息哈希值，以 0x 开头。                           |
| blockTime             | String       | 区块打包时间。                                        |
| blockSize             | Int          | 当前区块大小。                                        |
| proposerAddress       | Address      | 提案人地址。                                          |
| txs [{}]              | Object Array | 交易列表。                                            |
| txHash                | HexString    | 交易哈希值，以 0x 开头。                              |
| txTime                | String       | 交易时间。                                            |
| code                  | Uint32       | 交易结果码，200表示交易成功，其它值表示失败。           |
| log                   | String       | 交易结果描述。                                        |
| blockHash             | HexString    | 交易所在区块哈希值，以 0x 开头。                       |
| blockHeight           | Int64        | 交易所在区块高度。                                    |
| from                  | Address      | 交易签名人地址。                                      |
| nonce                 | Uint64       | 交易签名人交易计数值。                                 |
| gasLimit              | Uint64       | 最大燃料数量。                                        |
| fee                   | Uint64       | 交易手续费（单位cong）。                               |
| note                  | string       | 备注。                                                |
| messages [{}]         | Object Array | 消息列表。                                            |
| smcAddress            | Address      | 合约地址。                                            |
| smcName               | String       | 合约名称。                                            |
| method                | String       | 方法原型。                                            |
| to                    | Address      | 转账目的账户地址，仅当交易是BRC20标准转账时有效。       |
| value                 | string       | 转账金额（单位cong），仅当交易是BRC20标准转账时有效     |
| transferReceipts [{}] | Object Array | 收据列表。                                            |
| token                 | Address      | 代币地址                                              |
| from                  | Address      | 交易签名人地址。                                      |
| to                    | Address      | 转账目的账户地址，仅当交易是BRC20标准转账时有效。       |
| value                 | bn.Number    | 转账金额（单位cong），仅当交易是BRC20标准转账时有效。   |
| note                  | string       | 收据备注                                              |

#### 3.2.7 nonce

- **command**

  ```bash
  hdWallet nonce --address bcbAkTD... [--url https://...]
  ```

- **Input Parameters**

  | **语法** | **类型** | **注释**                                                     |
  | -------- | ------- | ------------------------------------------------------------ |
  | address  | Address | 账户地址。                                                   |
  | url      | String  | 钱包服务地址，可选项，默认调用本地服务。                       |

- **Output SUCCESS Example**

  ```json
  {
    "nonce": 5000
  }
  ```

- **Output SUCCESS Result**

  | **语法** | **类型** | **注释**                                                      |
  | -------- | ------ | ------------------------------------------------------------ |
  | nonce    | Uint64 | 指定地址在区块链上可用的下一个交易计数值。                   |

#### 3.2.8 allBalance

- **command**

  ```bash
  hdWallet allBalance --address bcbAkTD... [--url https://...]
  ```

- **Input Parameters**

  | **语法** | **类型** | **注释**                                                     |
  | -------- | ------- | ------------------------------------------------------------ |
  | address  | Address | 账户地址。                                                   |
  | url      | String  | 钱包服务地址，可选项，默认调用本地服务。                       |

- **Output SUCCESS Example**

  ```json
  [
     {
        "tokenAddress": "bcbALw9SqmqUWVjkB1bJUQyCKfnxDPhuN5Ej"，
        "tokenName": "bcb"，
        "balance": "2500000000"
     },
     {
        "tokenAddress": "bcbLVgb3odTfKC9Y9GeFnNWL9wmR4pwWiqwe"，
        "tokenName": "XT"，
        "balance": "10000000",
     }
  ]
  ```

- **Output SUCCESS Result**

  | **语法**     | **类型** | **注释**                                                     |
  | ------------ | :------: | ------------------------------------------------------------ |
  | tokenAddress | Address  | 代币地址。                                                   |
  | tokenName    |  String  | 代币名称。                                                   |
  | balance      |  String  | 账户余额（单位：Cong）。                                     |

#### 3.2.9 balance

- **command**

  ```bash
  hdWallet balance --address bcbAkTD... [--url https://...]
  ```

- **Input Parameters**

  | **语法** | **类型** | **注释**                                                      |
  | -------- | ------- | ------------------------------------------------------------ |
  | address  | Address | 账户地址。                                                   |
  | url      | String  | 钱包服务地址，可选项，默认调用本地服务。                       |

- **Output SUCCESS Example**

  ```json
  {
     "balance": "2500000000"
  }
  ```

- **Output SUCCESS Result**

  | **语法** | **类型** | **注释**                                                     |
  | -------- | ------ | ------------------------------------------------------------ |
  | balance  | String | 账户余额（单位：Cong）。                                     |

#### 3.2.10 balanceOfToken

- **command**

  ```bash
  hdWallet balanceOfToken --address bcbAkTD... --tokenAddress bcbCsRX... --tokenName XT [--url https://...]
  ```

- **Input Parameters**

  | **语法**     | **类型** | **注释**                                                      |
  | ------------ | ------- | ------------------------------------------------------------ |
  | address      | Address | 账户地址。                                                   |
  | tokenAddress | Address | 代币地址，与代币名称可以二选一，两个都有时必须一致。            |
  | tokenName    | String  | 代币名称，与代币地址可以二选一，两个都有时必须一致。            |
  | url          | String  | 钱包服务地址，可选项，默认调用本地服务。                       |

- **Output SUCCESS Example**

  ```json
  {
     "balance": "2500000000"
  }
  ```

- **Output SUCCESS Result**

  | **语法** | **类型** | **注释**                                                     |
  | -------- | ------ | ------------------------------------------------------------ |
  | balance  | String | 账户余额（单位：Cong）。                                      |

#### 3.2.11  transaction

- **command**

  ```bash
  hdWallet transaction --txHash 0x4E456161... [--url https://...]
  ```

- **Input Parameters**

  | **语法** | **类型**  | **注释** |
  | -------- | :-------: | ------------------------------------------------------------ |
  | txHash   | HexString | 指定交易哈希值，以 0x 开头。                                 |
  | url      |  String   | 钱包服务地址，可选项，默认调用本地服务。                     |

- **Output SUCCESS Example**

  ```json
  {
    "txHash": "0x4E456161A6580A1D34D86F1560DCFE6034F5E08FA31D7DCEBCCCC72A0DC73654",
    "txTime": "2018-12-27T14:26:19.251820644Z",
    "code": 200,
    "log": "Deliver tx succeed"
    "blockHash": "0x583E820E58D2FD00B1A7D66CDBB6B7C26B207925",
    "blockHeight": 2495461,
    "from": "bcbAkTDzHLf5udamub38GdepKe7nek66EHqY",
    "nonce": 117510,
    "gasLimit": 2500,
    "fee": 1500000,
    "note":"hello",
    "messages": [
      {
        "smcAddress": "bcbCsRXXMGkUJ8wRnrBUD7mQsMST4d53JRKJ",
        "smcName": "token-basic",
        "method": "Transfer(smc.Address,big.Int)smc.Error",
        "to": "bcbKuqW1qdsnD7mRsRooXMEkCBj2s9GLF9pn",
        "value": "683000000000"
      }
    ]
  }
  ```

- **Output SUCCESS Result**

  | **语法**                           |   **类型**   | **注释**                                                     |
  | ---------------------------------- | ------------ | ------------------------------------------------------------ |
  | txHash                             | HexString    | 交易哈希值，以 0x 开头。                                      |
  | txTime                             | String       | 交易时间。                                                    |
  | code                               | Uint32       | 交易结果码，200表示交易成功，其它值表示失败。                   |
  | log                                | String       | 交易结果描述。                                                |
  | blockHash                          | HexString    | 交易所在区块哈希值，以 0x 开头。                               |
  | blockHeight                        | Int64        | 交易所在区块高度。                                            |
  | from                               | Address      | 交易签名人地址。                                              |
  | nonce                              | Uint64       | 交易签名人交易计数值。                                        |
  | gasLimit                           | Uint64       | 最大燃料数量。                                                |
  | fee                                | Uint64       | 交易手续费（单位cong）。                                      |
  | note                               | String       | 备注。                                                       |
  | messages [{}]                      | Object Array | 消息列表。                                                   |
  | smcAddress                         | Address      | 合约地址。                                                   |
  | smcName                            | String       | 合约名称。                                                   |
  | method                             | String       | 方法原型。                                                   |
  | to                                 | Address      | 转账目的账户地址，仅当交易是BRC20标准转账时有效。              |
  | value                              | String       | 转账金额（单位cong），仅当交易是BRC20标准转账时有效。          |
