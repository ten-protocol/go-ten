# Open the log file in read mode
with open('/Users/WillHester/Desktop/Dev/_ten/go-ten/integration/.build/simulations/sim-log-2024-05-31_15-23-56-full-network-2007480794.txt', 'r') as logfile:
    # Open a new file to write the filtered lines
    with open('error_logfile.txt', 'w') as errorfile:
        # Loop through each line in the log file
        for line in logfile:
            # Check if the line starts with "[ERROR]"
            if line.startswith("ERROR["):
                # Write the line to the new file
                errorfile.write(line)
