import java.io.*;
import java.net.Socket;
import java.util.Scanner;

/**
 * Creates the Client to send messages to the proxy sever via TCP
 */
public class ChatClient {

    private final String host;
    private final int port;
    private Socket socket;
    private PrintWriter out;
    private BufferedReader in;
    private Scanner scanner;

    /**
     * Instantiates a ChatClient Object
     * @param host Host IP
     * @param port Port number
     */
    public ChatClient(String host, int port) {
        this.host = host;
        this.port = port;
    }

    /**
     * Start the client
     */
    public void start() {
        try {
            socket = new Socket(host, port);
            out = new PrintWriter(new OutputStreamWriter(socket.getOutputStream()), true);
            in = new BufferedReader(new InputStreamReader(socket.getInputStream()));
            scanner = new Scanner(System.in);

            // Start the thread to listen for server messages
            Thread serverListener = new Thread(new ServerHandler());
            serverListener.start();

            // Handle user input in the main thread
            handleUserInput();

            cleanup();

        } catch (IOException e) {
            System.err.println("I/O error: " + e.getMessage());
            cleanup();
        }
    }

    /**
     * Handle user input
     */
    private void handleUserInput() {
        String input;

        while (true) {
            System.out.print("> ");

            if (!scanner.hasNextLine()) {
                System.out.println("\nTerminating the client, EOF key detected.");
                break;
            }

            input = scanner.nextLine();

            // Send user input to the server
            out.println(input);
        }
    }

    /**
     * Method to clean up resources
     */
    private void cleanup() {
        try {
            if (socket != null && !socket.isClosed()) socket.close();
            if (scanner != null) scanner.close();
        } catch (IOException e) {
            System.err.println("Error during cleanup: " + e.getMessage());
        }
    }

    /**
     * Inner class to handle server messages
     */
    private class ServerHandler implements Runnable {
        @Override
        public void run() {
            String serverMessage;

            try {
                while ((serverMessage = in.readLine()) != null) {
                    System.out.println(serverMessage);
                    System.out.print("> ");
                }
                System.out.println("\nServer has closed the connection.");
            } catch (IOException e) {
                System.err.println("Connection to server lost: " + e.getMessage());
            } finally {
                cleanup();
            }
        }
    }

    public static void main(String[] args) {
        String host = "localhost";
        int port = 6666;

        ChatClient client = new ChatClient(host, port);
        client.start();
    }
}
