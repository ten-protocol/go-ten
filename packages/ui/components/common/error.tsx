import React from "react";
import NextErrorComponent from "next/error";
import Link from "next/link";
import { ErrorType } from "../../lib/interfaces/ui";

function ErrorMessage({
  statusText,
  message,
  showMessage,
  showStatusText,
}: any) {
  return (
    <div className="error-message">
      {showStatusText && <h3>{statusText}</h3>}
      {message && showMessage && (
        <p className="text-muted-foreground">{message}</p>
      )}
    </div>
  );
}

export function CustomError({
  showRedirectText = true,
  heading = "Oops! Something went wrong.",
  statusText = "500",
  message = "We're experiencing technical difficulties. Please try again later.",
  redirectText = "Home Page",
  isFullWidth,
  err,
  showMessage = true,
  showStatusText = true,
  statusCode,
  isModal,
  redirectLink = "/",
  children,
  ...props
}: ErrorType) {
  return (
    <section
      className="h-full flex flex-col justify-center items-center"
      {...props}
    >
      <main className={isFullWidth ? "max-w-full" : ""}>
        <div className="text-center">
          <h1 className="text-4xl font-extrabold mb-6">{heading}</h1>
          <div className={isFullWidth ? "w-full" : ""}>
            <ErrorMessage
              showStatusText={showStatusText}
              showMessage={showMessage}
              message={message}
              statusText={statusText}
            />
          </div>
          {showRedirectText && (
            <div>
              Go to{" "}
              <Link
                href={redirectLink}
                passHref
                className="text-primary pointer underline"
              >
                {redirectText}
              </Link>
            </div>
          )}
          {children}
        </div>
      </main>
    </section>
  );
}

// server-side error handling logic
CustomError.getInitialProps = async ({ res, err }: any) => {
  const statusCode = res ? res.statusCode : err?.statusCode || 404;
  const errorInitialProps = await NextErrorComponent.getInitialProps({
    res,
    err,
  } as any);
  errorInitialProps.statusCode = statusCode;

  // custom server-side error handling logic
  return statusCode < 500
    ? errorInitialProps
    : { ...errorInitialProps, statusCode };
};

export default CustomError;
