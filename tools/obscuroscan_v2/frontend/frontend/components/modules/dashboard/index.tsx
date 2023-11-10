import React from 'react'
import { CalendarDateRangePicker } from '@/components/date-range-picker'
import { CardHeader, CardTitle, CardContent, Card } from '@/components/ui/card'
import { TabsList, TabsTrigger, TabsContent, Tabs } from '@/components/ui/tabs'
import { LayersIcon, FileTextIcon, ReaderIcon, CubeIcon, RocketIcon } from '@radix-ui/react-icons'

import { RecentBatches } from './recent-batches'
import { RecentRollups } from './recent-rollups'
import { RecentTransactions } from './recent-transactions'
import { Button } from '@/components/ui/button'

const DASHBOARD_DATA = [
  {
    title: 'Ether Price',
    value: '$1967.89',
    change: '+20.1% from last month',
    icon: RocketIcon
  },
  {
    title: 'Latest Batch',
    value: '2061',
    change: '+20.1% from last month',
    icon: LayersIcon
  },
  {
    title: 'Latest Rollup',
    value: '0xbaa7c0288013169ca1dc9115b5ce496a31a5972f16b6975a62cc09cc740f294e',
    change: '+20.1% from last month',
    icon: CubeIcon
  },
  {
    title: 'Transactions',
    value: '5',
    change: '+20.1% from last month',
    icon: ReaderIcon
  },
  {
    title: 'Contracts',
    value: '3',
    change: '+20.1% from last month',
    icon: FileTextIcon
  }
]

export default function Dashboard() {
  return (
    <div className="flex-1 space-y-4 p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Obscuroscan</h2>
        <div className="flex items-center space-x-2">
          <CalendarDateRangePicker />
          <Button>Download</Button>
        </div>
      </div>
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="analytics" disabled>
            Analytics
          </TabsTrigger>
          <TabsTrigger value="reports" disabled>
            Reports
          </TabsTrigger>
          <TabsTrigger value="notifications" disabled>
            Notifications
          </TabsTrigger>
        </TabsList>
        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {DASHBOARD_DATA.map((item: any, index) => (
              <Card key={index}>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">{item.title}</CardTitle>
                  {React.createElement(item.icon)}
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold truncate">{item.value}</div>
                  <p className="text-xs text-muted-foreground">{item.change}</p>
                </CardContent>
              </Card>
            ))}
          </div>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-9">
            <Card className="col-span-3">
              <CardHeader>
                <CardTitle>Recent Batches</CardTitle>
              </CardHeader>
              <CardContent>
                <RecentBatches />
              </CardContent>
            </Card>
            <Card className="col-span-3">
              <CardHeader>
                <CardTitle>Recent Rollups</CardTitle>
              </CardHeader>
              <CardContent>
                <RecentRollups />
              </CardContent>
            </Card>
            <Card className="col-span-3">
              <CardHeader>
                <CardTitle>Recent Transactions</CardTitle>
              </CardHeader>
              <CardContent>
                <RecentTransactions />
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
