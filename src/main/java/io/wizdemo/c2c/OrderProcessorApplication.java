package io.wizdemo.c2c;

import com.fasterxml.jackson.databind.ObjectMapper;
import java.util.LinkedHashMap;
import java.util.Map;
import org.apache.commons.text.StringSubstitutor;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.yaml.snakeyaml.Yaml;

@SpringBootApplication
@RestController
public class OrderProcessorApplication {
  private static final Logger LOGGER = LogManager.getLogger(OrderProcessorApplication.class);
  private final ObjectMapper objectMapper = new ObjectMapper();

  @GetMapping("/healthz")
  public Map<String, Object> healthz(
      @RequestParam(name = "template", defaultValue = "order-processor-${java:version}")
          String template) throws Exception {
    LOGGER.info("Health check requested for c2c-hosted-applications");
    StringSubstitutor interpolator = StringSubstitutor.createInterpolator();
    String rendered = interpolator.replace(template);

    Yaml yaml = new Yaml();
    Map<String, Object> parsed = yaml.load("service: order-processor\nversion: 1.0.0\n");

    Map<String, Object> response = new LinkedHashMap<>();
    response.put("status", "ok");
    response.put("service", parsed.get("service"));
    response.put("version", parsed.get("version"));
    response.put("rendered", rendered);
    response.put("serialized", objectMapper.writeValueAsString(parsed));
    return response;
  }

  public static void main(String[] args) {
    SpringApplication.run(OrderProcessorApplication.class, args);
  }
}
