import React from 'react'
import { Metadata } from 'next'
import { LayersIcon, FileTextIcon, ReaderIcon, CubeIcon, RocketIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { CalendarDateRangePicker } from '@/components/date-range-picker'
import { MainNav } from '@/components/main-nav'
import { Search } from '@/components/search'
import TeamSwitcher from '@/components/team-switcher'
import { ModeToggle } from '@/components/mode-toggle'
import { RecentBatches } from '@/components/modules/dashboard/recent-batches'
import { RecentTransactions } from '@/components/modules/dashboard/recent-transactions'
export const metadata: Metadata = {
  title: 'Dashboard',
  description: 'Example dashboard app built using the components.'
}

export default function DashboardPage() {
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

  return (
    <>
      <div className="hidden flex-col md:flex">
        <div className="border-b">
          <div className="flex h-16 items-center px-4">
            <TeamSwitcher />
            <MainNav className="mx-6" />
            <div className="ml-auto flex items-center space-x-4">
              <Search />
              <ModeToggle />
            </div>
          </div>
        </div>
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
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
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
      </div>
    </>
  )
}
