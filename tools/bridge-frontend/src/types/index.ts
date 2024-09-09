import React from "react";

export interface SeoProps {
  title: string;
  description: string;
  canonicalUrl: string;
  ogTwitterImage: string;
  ogImageUrl: string;
  ogType: string;
  includeDefaultKeywords?: boolean;
  children?: React.ReactNode;
}
