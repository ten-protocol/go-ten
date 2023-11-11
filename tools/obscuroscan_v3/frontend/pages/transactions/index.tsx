import React, { useEffect, useState } from 'react'
import { Payment, columns } from '@/components/modules/personal/transactions/columns'
import { DataTable } from '@/components/modules/common/data-table'
import Layout from '@/components/layouts/default-layout'

async function getData(): Promise<Payment[]> {
  return [
    {
      id: '728ed52f',
      amount: 100,
      status: 'pending',
      email: 'm@example.com'
    },
    {
      id: '489e1d42',
      amount: 125,
      status: 'processing',
      email: 'example@gmail.com'
    }
  ]
}

function Transactions() {
  const [data, setData] = useState<Payment[]>([])

  useEffect(() => {
    const fetchData = async () => {
      try {
        const result = await getData()
        setData(result)
      } catch (error) {
        console.error('Error fetching data:', error)
      }
    }

    fetchData()
  }, [])

  return (
    <Layout>
      <div className="container mx-auto py-10">
        <div className="flex items-center justify-between space-y-2">
          <h2 className="text-3xl font-bold tracking-tight mb-4">Transactions</h2>
        </div>
        {data.length > 0 ? <DataTable columns={columns} data={data} /> : <p>Loading...</p>}
      </div>
    </Layout>
  )
}

export default Transactions
