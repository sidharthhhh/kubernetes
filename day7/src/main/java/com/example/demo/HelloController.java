package com.example.demo;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import java.net.InetAddress;
import java.net.UnknownHostException;

@RestController
public class HelloController {

	@GetMapping("/hello")
	public String hello() {
		String hostname = "Unknown";
		try {
			hostname = InetAddress.getLocalHost().getHostName();
		} catch (UnknownHostException e) {
			e.printStackTrace();
		}
		return "Hello from Spring Boot on Kubernetes! Pod: " + hostname;
	}

    @GetMapping("/health")
    public String health() {
        return "OK";
    }
}
