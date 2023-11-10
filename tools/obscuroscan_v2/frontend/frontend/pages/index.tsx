import React from 'react'
import { Metadata } from 'next'
// import { useTransactions } from '@/components/hooks/useTransactions'
import Layout from '@/components/layouts/default-layout'
import Dashboard from '@/components/modules/dashboard'

export const metadata: Metadata = {
  title: 'Dashboard',
  description: 'Example dashboard app built using the components.'
}

export default function DashboardPage() {
  // const { transactions, isTransactionsLoading, transactionCount, isTransactionCountLoading } =
  //   useTransactions()
  // console.log('🚀 ~ file: index.tsx:23 ~ DashboardPage ~ transactions:', transactions)

  return (
    <Layout>
      <Dashboard />
    </Layout>
  )
}
