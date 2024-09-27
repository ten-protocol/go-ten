import { IWalletState } from "../interfaces/wallet";

export type StoreSet = (
  partial:
    | IWalletState
    | Partial<IWalletState>
    | ((state: IWalletState) => IWalletState | Partial<IWalletState>),
  replace?: boolean | undefined
) => void;

export type StoreGet = () => IWalletState;

export interface ResponseDataInterface<T> {
  result: T;
  errors?: string[] | string;
  item: T;

  message: string;
  pagination?: PaginationInterface;
  success: string;
}

export interface PaginationInterface {
  page: number;
  perPage: number;
  total: number;
  totalPages: number;
}

export interface DocumentContentInterface {
  heading: string;
  content: string[];
}

export interface DocumentInterface {
  title: string;
  subHeading: string;
  content: DocumentContentInterface[];
}
