import { ethers } from "ethers";

export interface SeoProps {
  title: string;
  description: string;
  canonicalUrl: string;
  ogTwitterImage: string;
  ogImageUrl: string;
  ogType: string;
  children: React.ReactNode;
}

export interface ErrorType {
  statusCode?: number;
  showRedirectText?: boolean;
  heading?: string;
  statusText?: string;
  message?: string;
  redirectText?: string;
  customPageTitle?: string;
  isFullWidth?: boolean;
  style?: React.CSSProperties;
  hasGetInitialPropsRun?: boolean;
  err?: Error;
  showMessage?: boolean;
  showStatusText?: boolean;
  isModal?: boolean;
  redirectLink?: string;
  children?: React.ReactNode;
  [key: string]: any;
}

export interface IconProps {
  width?: string;
  height?: string;
  fill?: string;
  stroke?: string;
  strokeWidth?: string;
  className?: string;
  isActive?: boolean;
  onClick?: () => void;
}

export interface GetInfinitePagesInterface<T> {
  nextId?: number;
  previousId?: number;
  data: T;
  count: number;
}

export interface PaginationInterface {
  page: number;
  perPage: number;
  total: number;
  totalPages: number;
}

export interface ResponseDataInterface<T> {
  result: T;
  item: T;
  message: string;
  pagination?: PaginationInterface;
  success: boolean;
}

export type NavLink = {
  label: string;
  href?: string;
  isDropdown?: boolean;
  isExternal?: boolean;
  subNavLinks?: NavLink[];
  icon?: React.ElementType;
};

export enum ToastType {
  INFO = "info",
  SUCCESS = "success",
  WARNING = "warning",
  DESTRUCTIVE = "destructive",
  DEFAULT = "default",
}

export interface SeoProps {
  title: string;
  description: string;
  canonicalUrl: string;
  ogTwitterImage: string;
  ogImageUrl: string;
  ogType: string;
  children: React.ReactNode;
}

export interface IconProps {
  width?: string;
  height?: string;
  fill?: string;
  stroke?: string;
  strokeWidth?: string;
  className?: string;
  isActive?: boolean;
  onClick?: () => void;
}

export interface GetInfinitePagesInterface<T> {
  nextId?: number;
  previousId?: number;
  data: T;
  count: number;
}

export interface PaginationInterface {
  page: number;
  perPage: number;
  total: number;
  totalPages: number;
}

export interface ResponseDataInterface<T> {
  result: T;
  item: T;
  message: string;
  pagination?: PaginationInterface;
  success: boolean;
}

export interface WalletConnectionContextType {}

export interface Props {
  children: React.ReactNode;
}

export interface State {
  hasError: boolean;
}

export interface WalletConnectionProviderProps {
  children: React.ReactNode;
}

export enum L1Network {
  MAINNET = 1,
  SEPOLIA = 11155111,
  UAT = 1337,
  DEV = 1337,
  LOCAL = 1337,
}

export enum L2Network {
  SEPOLIA = 443,
  UAT = 443,
  DEV = 443,
  LOCAL = 443,
}

export type Environment = "uat-testnet" | "sepolia-testnet" | "dev-testnet";

export interface IChain {
  name: string;
  value: string;
  isNative: boolean;
  isEnabled: boolean;
  chainId?: string;
}

export interface IToken {
  name: string;
  value: string;
  isNative: boolean;
  isEnabled: boolean;
  address: string;
}

export interface INetworkConfig {
  ManagementContractAddress: string;
  L1StartHash: string;
  MessageBusAddress: string;
  L2MessageBusAddress: string;
  ImportantContracts: {
    L1Bridge: string;
    L1CrossChainMessenger: string;
    L2Bridge: string;
    L2CrossChainMessenger: string;
  };
}

export interface IContractState {
  bridgeContract: ethers.Contract | null;
  managementContract: ethers.Contract | null;
  messageBusContract: ethers.Contract | null;
  messageBusAddress: string;
  bridgeAddress: string;
  setContractState: (state: Partial<IContractState>) => void;
}

export interface IWalletState {
  provider: ethers.providers.Web3Provider | null;
  signer: ethers.Signer | null;
  address: string;
  walletConnected: boolean;
  isL1ToL2: boolean;
  isWrongNetwork: boolean;
  loading: boolean;
  initializeProvider: () => void;
  connectWallet: () => void;
  disconnectWallet: () => void;
  switchNetwork: () => void;
  restoreWalletState: () => void;
}

export type ButtonVariants =
  | "outline"
  | "link"
  | "default"
  | "destructive"
  | "secondary"
  | "ghost"
  | "clear";

export enum ItemPosition {
  FIRST = "first",
  LAST = "last",
}

export enum TransactionStatus {
  Success = "Success",
  Failure = "Failure",
  Pending = "Pending",
}

export type Transactions = {
  blockNumber: number;
  blockHash: string;
  transactionIndex: number;
  removed: boolean;
  address: string;
  data: string;
  topics: string[];
  transactionHash: string;
  logIndex: number;
  status: TransactionStatus;
};

export type StoreSet = (
  partial:
    | IWalletState
    | Partial<IWalletState>
    | ((state: IWalletState) => IWalletState | Partial<IWalletState>),
  replace?: boolean | undefined
) => void;

export type StoreGet = () => IWalletState;

export enum TransactionStep {
  TransactionSubmission = "transactionSubmission",
  TransactionConfirmation = "transactionConfirmation",
  EventDataExtraction = "eventDataExtraction",
  MerkleTreeConstruction = "merkleTreeConstruction",
  GasEstimation = "gasEstimation",
  RelaySubmission = "relaySubmission",
}

export enum IErrorMessages {
  UnknownAccount = "unknown account",
  InsufficientFunds = "insufficient funds",
  UserDeniedTransactionSignature = "User denied transaction signature",
  UserRejectedTheRequest = "User rejected the request",
  ExecutionReverted = "execution reverted",
  RateLimitExceeded = "rate limit exceeded",
}

export interface IPendingTx {
  receiver?: string;
  value?: string;
  token?: string;
  txHash?: string;
  resumeStep?: TransactionStep;
  txResponse?: ethers.providers.TransactionResponse;
  txReceipt?: ethers.providers.TransactionReceipt;
  valueTransferEventData?: ethers.utils.LogDescription;
  block?: any;
  tree?: any;
  proof?: any;
  gasLimit?: ethers.BigNumber | null;
  timestamp?: number;
}
