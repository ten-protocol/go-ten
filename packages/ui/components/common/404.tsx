import { ErrorType } from "../../lib/interfaces/ui";
import CustomError from "./error";

export function Custom404Error({
  customPageTitle,
  showRedirectText,
  redirectText,
  isFullWidth,
  message,
  showMessage = true,
  redirectLink,
  children,
}: ErrorType) {
  return (
    <CustomError
      heading={` ${customPageTitle || "Oops! Page"} Not Found`}
      statusText={`We can't seem to find the ${
        customPageTitle || "page"
      } you're looking for.`}
      statusCode={404}
      showRedirectText={showRedirectText}
      redirectText={redirectText || "Home Page"}
      message={
        message ||
        `The ${
          customPageTitle || "page"
        } you are looking for might have been removed, had its name changed, or is temporarily unavailable.`
      }
      isFullWidth={isFullWidth}
      showMessage={showMessage}
      redirectLink={redirectLink}
    >
      {children}
    </CustomError>
  );
}

export default Custom404Error;
