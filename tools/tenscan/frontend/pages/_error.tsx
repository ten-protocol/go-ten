import CustomError from "@repo/ui/components/common/error";
import React from "react";

export default function CustomErrorPage({ statusCode, err }: any) {
  return (
    <CustomError
      statusCode={statusCode}
      err={err}
      heading={
        statusCode === 404 ? "Page Not Found" : "Oops! Something went wrong."
      }
      message={
        statusCode === 404
          ? "Sorry, the page you're looking for doesn't exist."
          : "We're experiencing technical difficulties. Please try again later."
      }
      redirectLink="/"
      redirectText="Back to Home"
    />
  );
}

// custom server-side props
CustomErrorPage.getInitialProps = async ({ res, err }: any) => {
  const statusCode = res ? res.statusCode : err?.statusCode || 404;
  return { statusCode, err };
};
