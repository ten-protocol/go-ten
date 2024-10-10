import { ErrorType } from "../../lib/interfaces/ui";
import CustomError from "./error";

function Custom500Error({
  customPageTitle,
  message,
  showRedirectText,
  redirectText,
  err,
  redirectLink,
  children,
}: ErrorType) {
  return (
    <CustomError
      heading={"Oops! Something went wrong."}
      message={
        message ||
        "We're experiencing technical difficulties. Please try again later."
      }
      statusText={customPageTitle || `An Error occured`}
      statusCode={500}
      showRedirectText={showRedirectText || true}
      redirectText={redirectText || "Home Page"}
      err={err}
      redirectLink={redirectLink}
    >
      {children}
    </CustomError>
  );
}

export default Custom500Error;
