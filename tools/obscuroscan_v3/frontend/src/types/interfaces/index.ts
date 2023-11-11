import React from 'react'

export interface ErrorType {
  statusCode?: number
  showRedirectText?: boolean
  heading?: string
  statusText?: string
  message?: string
  redirectText?: string
  customPageTitle?: string
  isFullWidth?: boolean
  style?: React.CSSProperties
  hasGetInitialPropsRun?: boolean
  err?: Error
  showMessage?: boolean
  showStatusText?: boolean
  isModal?: boolean
  redirectLink?: string
  children?: React.ReactNode
  [key: string]: any
}

export interface IconProps {
  width?: string
  height?: string
  fill?: string
  stroke?: string
  strokeWidth?: string
  className?: string
  isActive?: boolean
  onClick?: () => void
}

export interface GetInfinitePagesInterface<T> {
  nextId?: number
  previousId?: number
  data: T
  count: number
}

export interface PaginationInterface {
  page: number
  perPage: number
  total: number
  totalPages: number
}

export interface ResponseDataInterface<T> {
  data: T
  message: string
  pagination?: PaginationInterface
  success: string
}
