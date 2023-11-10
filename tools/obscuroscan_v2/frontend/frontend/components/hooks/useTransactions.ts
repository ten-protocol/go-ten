import { getTransactions, getTransactionCount } from '@/api/transactions'
import { useQuery } from '@tanstack/react-query'

export const useTransactions = () => {
  const { data: transactions, isLoading: isTransactionsLoading } = useQuery({
    queryKey: ['transactions'],
    queryFn: () => getTransactions()
  })

  const { data: transactionCount, isLoading: isTransactionCountLoading } = useQuery({
    queryKey: ['transactionCount'],
    queryFn: () => getTransactionCount()
  })

  return { transactions, isTransactionsLoading, transactionCount, isTransactionCountLoading }
}
