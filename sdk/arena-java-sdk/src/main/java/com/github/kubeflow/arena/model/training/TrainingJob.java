package com.github.kubeflow.arena.model.training;

import com.github.kubeflow.arena.enums.TrainingJobType;

import java.util.List;

public class TrainingJob {
    private String name;
    private TrainingJobType jobType;
    private List<String> args;
    private String command;

    public TrainingJob(String name, TrainingJobType jobType, List<String> args, String command) {
        this.name = name;
        this.jobType = jobType;
        this.args = args;
        this.command = command;
    }

    public String name() {
        return this.name;
    }

    public TrainingJobType getType() {
        return this.jobType;
    }

    public List<String> getArgs() {
        return this.args;
    }

    public String getCommand() {
        if (this.command == null) {
            return "";
        }
        return this.command;
    }
}
